package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"errors"

	"bitbucket.org/SlothNinja/chat"
	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/sn"
	"bitbucket.org/SlothNinja/stack"
	"bitbucket.org/SlothNinja/status"
	"bitbucket.org/SlothNinja/store"
	"bitbucket.org/SlothNinja/user"
	"cloud.google.com/go/datastore"
	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
)

const param = "hid"

type action func(game, *gin.Context) (game, error)
type stackFunc func(stack.Stack) stack.Stack

type server struct {
	store.Store
}

var noServer = &server{}

func (s server) show() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		s, err := s.init(c)
		if err != nil {
			jerr(c, err)
			return
		}

		g, err := s.getGame(c, param)
		if err != nil {
			jerr(c, err)
			return
		}

		cu := user.Current(c)
		g = g.updateClickablesFor(cu)
		c.JSON(http.StatusOK, gin.H{gameKey: g})
	}
}

func (s server) newInvitation() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu := user.Current(c)
		if cu == user.None {
			jerr(c, errUserNotFound)
			return
		}

		e := newHeaderEntity(newGame(0))

		// Default Values
		e.Title = fmt.Sprintf("%s's %s", cu.Name, randomdata.SillyName())
		e.NumPlayers = 2
		e.TwoThiefVariant = false

		c.JSON(http.StatusOK, gin.H{"header": e})
	}
}

func (s server) createInvitation() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu := user.Current(c)
		if cu == user.None {
			jerr(c, errUserNotFound)
			return
		}

		s, err := s.init(c)
		if err != nil {
			jerr(c, err)
			return
		}

		type jH struct {
			Title      string `json:"title" binding:"min=4,max=30"`
			NumPlayers int    `json:"numPlayers" binding:"gte=2,lte=4"`
			TwoThief   bool   `json:"twoThief"`
			Password   string `json:"password"`
		}

		type jD struct {
			Header jH `json:"header"`
		}

		jData := new(jD)

		err = c.ShouldBindJSON(jData)
		if err != nil {
			jerr(c, err)
			return
		}

		t := time.Now()
		invitation := newHeaderEntity(newGame(0))
		invitation.Title = jData.Header.Title
		invitation.TwoThiefVariant = jData.Header.TwoThief
		invitation.NumPlayers = jData.Header.NumPlayers
		invitation.Password = jData.Header.Password
		invitation.Creator = cu
		invitation = invitation.AddUser(cu)
		invitation.Status = status.Recruiting
		invitation.CreatedAt, invitation.UpdatedAt = t, t

		_, err = s.RunInTransaction(c, func(tx *datastore.Transaction) error {
			_, err := tx.Put(invitation.Key, &invitation)
			if err != nil {
				return fmt.Errorf("unable to put header: %w", err)
			}

			m := chat.NewMLog(c, invitation.ID())
			m.CreatedAt, m.UpdatedAt = t, t
			_, err = tx.Put(m.Key, m)
			if err != nil {
				return fmt.Errorf("unable to put chat: %w", err)
			}
			return nil
		})

		if err != nil {
			jerr(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s created game %q", cu.Name, invitation.Title)})
	}
}

func (s server) acceptInvitation(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu := user.Current(c)
		if cu == user.None {
			jerr(c, errUserNotFound)
			return
		}

		hid, err := strconv.ParseInt(c.Param(param), 10, 64)
		if err != nil {
			jerr(c, fmt.Errorf("unable to get header id: %w", err))
			return
		}

		invitation := newHeaderEntity(newGame(hid))

		s, err := s.init(c)
		if err != nil {
			jerr(c, err)
			return
		}

		err = s.Get(c, invitation.Key, &invitation)
		if err != nil {
			jerr(c, fmt.Errorf("unable to get header: %w", err))
			return
		}

		invitation, start, err := invitation.Accept(c, cu)
		if err != nil {
			jerr(c, err)
			return
		}

		if start {
			invitation.Status = status.Starting

			g := newGame(invitation.ID())
			g.Header = invitation.Header
			g = g.start()

			_, err := s.RunInTransaction(c, func(tx *datastore.Transaction) error {
				t := time.Now()
				g.UpdatedAt, g.StartedAt = t, t
				_, err := tx.PutMulti(g.withHistory())
				return err
			})

			if err != nil {
				jerr(c, err)
				return
			}

			sendTurnNotification(c, g.Header, g.CPUserIndices[0])
			c.JSON(http.StatusOK, gin.H{
				"header":  invitation,
				"message": cu.Name + " joined game " + invitation.Title,
			})
			return
		}

		_, err = s.Put(c, invitation.Key, &invitation)
		if err != nil {
			jerr(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"header":  invitation,
			"message": cu.Name + " joined game " + invitation.Title,
		})
	}
}

func (s server) dropInvitation(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu := user.Current(c)
		if cu == user.None {
			jerr(c, errUserNotFound)
			return
		}

		s, err := s.init(c)
		if err != nil {
			jerr(c, err)
			return
		}

		hid, err := strconv.ParseInt(c.Param(param), 10, 64)
		if err != nil {
			jerr(c, fmt.Errorf("unable to get header id: %w", err))
			return
		}

		invitation := newHeaderEntity(newGame(hid))

		err = s.Get(c, invitation.Key, &invitation)
		if err != nil {
			jerr(c, err)
			return
		}

		invitation, err = invitation.Drop(cu)
		if err != nil {
			jerr(c, err)
			return
		}

		if len(invitation.Users) == 0 {
			invitation.Status = status.Aborted
		}

		_, err = s.Put(c, invitation.Key, &invitation)
		if err != nil {
			jerr(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"header":  invitation,
			"message": cu.Name + " left invitation for game \"" + invitation.Title + "\"",
		})
	}
}

func (s server) playCard() gin.HandlerFunc {
	return s.update(param, (game).PlayCard)
}

func (s server) selectThief() gin.HandlerFunc {
	return s.update(param, (game).SelectThief)
}

func (s server) pass() gin.HandlerFunc {
	return s.update(param, (game).Pass)
}

// var unexpectedErr = "Unexpected error.  Try again."

func jerr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errValidation):
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
	default:
		log.Debugf(err.Error())
		c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
	}
}

func (s server) update(param string, act action) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		s, err := s.init(c)
		if err != nil {
			jerr(c, err)
			return
		}

		g, err := s.getGame(c, param)
		if err != nil {
			jerr(c, err)
			return
		}

		g, err = act(g, c)
		if err != nil {
			jerr(c, err)
			return
		}

		_, err = s.RunInTransaction(c, func(tx *datastore.Transaction) error {
			_, err = tx.PutMulti(g.withHistory())
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			jerr(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{gameKey: g})
	}
}

func (s server) init(c *gin.Context) (server, error) {
	if s.Store != nil {
		return s, nil
	}

	var err error
	s.Store, err = datastore.NewClient(c, "")
	if err != nil {
		return noServer, err
	}
	return s, nil
}

func (s server) updateStack(param string, act action, animate bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, err := s.init(c)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		g, err := s.getGame(c, param)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		g, err = act(g, c)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		_, err = s.Put(c, g.Key, g)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		// Essentially reloads and presents game from datastore
		s.show()(c)
	}
}

func (s server) reset(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu := user.Current(c)
		if cu == user.None {
			jerr(c, errUserNotFound)
			return
		}

		s, err := s.init(c)
		if err != nil {
			jerr(c, err)
			return
		}

		g, err := s.getGame(c, param)
		if err != nil {
			jerr(c, err)
			return
		}

		cp := g.currentPlayerFor(cu)
		if cp.id == noPID || !g.CPorAdmin(cp.id, cu) {
			jerr(c, fmt.Errorf("only the current player can reset a move: %w", errValidation))
			return
		}

		var changed bool
		g.Stack, changed = g.Stack.Reset()
		if !changed {
			g = g.updateClickablesFor(cu)
			c.JSON(http.StatusOK, gin.H{gameKey: g})
			return
		}

		h, err := s.getHistory(c, g)
		if err != nil {
			jerr(c, err)
			return
		}

		g.Header, g.state = h.Header, h.state
		g.CreatedAt, g.UpdatedAt = h.CreatedAt, h.UpdatedAt
		_, err = s.RunInTransaction(c, func(tx *datastore.Transaction) error {
			_, err := tx.PutMulti(g.withoutHistory())
			return err
		})

		if err != nil {
			jerr(c, err)
			return
		}

		g = g.updateClickablesFor(cu)
		c.JSON(http.StatusOK, gin.H{gameKey: g})
	}
}

func (s server) undo(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu := user.Current(c)
		if cu == user.None {
			jerr(c, errUserNotFound)
			return
		}

		s, err := s.init(c)
		if err != nil {
			jerr(c, err)
			return
		}

		g, err := s.getGame(c, param)
		if err != nil {
			jerr(c, err)
			return
		}

		cp := g.currentPlayerFor(cu)
		if cp.id == noPID || !g.CPorAdmin(cp.id, cu) {
			jerr(c, fmt.Errorf("only the current player can reset a move: %w", errValidation))
			return
		}

		var changed bool
		g.Stack, changed = g.Stack.Undo()
		if !changed {
			g = g.updateClickablesFor(cu)
			c.JSON(http.StatusOK, gin.H{gameKey: g})
			return
		}

		h, err := s.getHistory(c, g)
		if err != nil {
			jerr(c, err)
			return
		}

		g.Header, g.state, g.CreatedAt, g.UpdatedAt = h.Header, h.state, h.CreatedAt, h.UpdatedAt
		_, err = s.RunInTransaction(c, func(tx *datastore.Transaction) error {
			log.Debugf("saved g.Stack: %#v", g.Stack)
			_, err := tx.PutMulti(g.withoutHistory())
			return err
		})

		if err != nil {
			jerr(c, err)
			return
		}

		g = g.updateClickablesFor(cu)
		c.JSON(http.StatusOK, gin.H{gameKey: g})
	}
}

func (s server) redo(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu := user.Current(c)
		if cu == user.None {
			jerr(c, errUserNotFound)
			return
		}

		s, err := s.init(c)
		if err != nil {
			jerr(c, err)
			return
		}

		g, err := s.getGame(c, param)
		if err != nil {
			jerr(c, err)
			return
		}

		cp := g.currentPlayerFor(cu)
		if !g.CPorAdmin(cp.id, cu) {
			jerr(c, fmt.Errorf("only the current player can reset a move: %w", errValidation))
			return
		}

		var changed bool
		g.Stack, changed = g.Stack.Redo()
		if !changed {
			g = g.updateClickablesFor(cu)
			c.JSON(http.StatusOK, gin.H{gameKey: g})
			return
		}

		h, err := s.getHistory(c, g)
		if err != nil {
			jerr(c, err)
			return
		}

		g.Header, g.state, g.CreatedAt, g.UpdatedAt = h.Header, h.state, h.CreatedAt, h.UpdatedAt
		_, err = s.RunInTransaction(c, func(tx *datastore.Transaction) error {
			_, err := tx.PutMulti(g.withoutHistory())
			return err
		})

		if err != nil {
			jerr(c, err)
			return
		}

		g = g.updateClickablesFor(cu)
		c.JSON(http.StatusOK, gin.H{gameKey: g})
	}
}

// func (s *server) before(c *gin.Context, param string) (*game, *user.User2, error) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	g, err := s.getGame(c, param)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	cu := user.Current(c)
// 	return g, cu, nil
// }

const (
	playAnimations = true
	skipAnimations = false
)

func (s server) getParam(c *gin.Context, param string) (int64, error) {
	i, err := sn.Int64Param(c, param)
	if err != nil {
		return 0, fmt.Errorf("unable to get param: %w", errValidation)
	}
	return i, nil
}

// func (s server) getHeader(c *gin.Context, param string) (*gotHeader, error) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	hid, err := s.getParam(c, param)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	h := New(hid)
//
// 	err = s.Get(c, h.Key, h)
// 	return h, err
// }

func (s server) getGame(c *gin.Context, param string) (game, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	id, err := s.getParam(c, param)
	if err != nil {
		return noGame, err
	}

	g := newGame(id)
	err = s.Get(c, g.Key, &g)
	return g, err
}

func (s server) getHistory(c *gin.Context, g game) (history, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	h := newHistory(g)
	err := s.Get(c, h.Key, &h)
	return h, err
}

// func (s server) getState(c *gin.Context, h gotHeader) (state.state, error) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	st := state.New(h.Current, h.Key)
// 	err := s.Get(c, st.Key, st)
// 	return st, err
// }

// func returnJSON(c *gin.Context, g game, msg string) {
// 	c.JSON(http.StatusOK, struct {
// 		game    game  `json:"game"`
// 		Message string `json:"message"`
// 	}{
// 		game:    g,
// 		Message: msg,
// 	})
// }
