package main

import (
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/sn"
	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
)

func (s server) finish(param, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf(msgEnter)
		defer log.Debugf(msgExit)

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

		g, err = g.FinishTurn(c)
		if err != nil {
			jerr(c, err)
			return
		}

		// var rs []result

		if g.Phase == phaseGameOver {
			g.Status = sn.Ending
			_, err = g.endgame(c)
			// rs, err = g.endgame(c)
			if err != nil {
				jerr(c, err)
				return
			}
		}

		t := time.Now()
		g.UpdatedAt = t

		_, err = s.RunInTransaction(c, func(tx *datastore.Transaction) error {
			_, err = tx.PutMulti(g.withHistory())
			if err != nil {
				return err
			}
			// err = sendEndGameNotifications(c, g.Header, rs)
			// if err != nil {
			// 	log.Warningf(err.Error())
			// }
			return nil
		})

		if err != nil {
			jerr(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{gameKey: g})
	}
}

// FinishTurn provides a handler for finishing a turn.
func (g game) FinishTurn(c *gin.Context) (game, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	switch g.Phase {
	case phasePlaceThieves:
		return g.placeThievesFinishTurn(c)
	case phaseClaimItem:
		return g.moveThiefFinishTurn(c)
	default:
		return g, errWrongPhase
	}
}

func (g game) validateFinishTurn(c *gin.Context) (player, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cu := user.Current(c)
	if cu == user.None {
		return noPlayer, errUserNotFound
	}

	cp := g.currentPlayerFor(cu)
	switch {
	case !g.CPorAdmin(cp.id, cu):
		return noPlayer, errNotCPorAdmin
	case !cp.performedAction:
		return noPlayer, errActionNotPerformed
	default:
		return cp, nil
	}
}

func (g game) nextPlayer(cp player) (player, bool) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	index, found := playerFindIndex(cp, g.players)
	if !found {
		return noPlayer, false
	}
	return playerByIndex(index+1, g.players), true
}

func (g game) previousPlayer(cp player) (player, bool) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)
	index, found := playerFindIndex(cp, g.players)
	if !found {
		return noPlayer, false
	}
	return playerByIndex(index-1, g.players), true
}

func (g game) placeThievesFinishTurn(c *gin.Context) (game, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, err := g.validatePlaceThievesFinishTurn(c)
	if err != nil {
		return g, err
	}

	np, found := g.previousPlayer(cp)
	if !found {
		return g, errPlayerNotFound
	}

	newTurn := cp.Equal(g.players[0])
	if newTurn {
		g.Turn++
	}

	numThieves := 3
	if g.TwoThiefVariant {
		numThieves = 2
	}

	if newTurn && g.Turn > numThieves {
		g = g.startCardPlay()
		np = cp
	}

	g.CPUserIndices = []int{np.id}

	g.Stack = g.Commit()
	g = g.beginningOfTurnReset(np)

	cu := user.Current(c)
	g = g.updateClickablesFor(cu)

	return g, nil
}

func (g game) validatePlaceThievesFinishTurn(c *gin.Context) (player, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, err := g.validateFinishTurn(c)
	switch {
	case err != nil:
		return noPlayer, err
	case g.Phase != phasePlaceThieves:
		return noPlayer, errWrongPhase
	default:
		return cp, nil
	}
}

func (g game) moveThiefNextPlayer(p player) (player, bool) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	for !allPassed(g.players) {
		np, found := g.nextPlayer(p)
		switch {
		case !found:
			return noPlayer, false
		case !np.passed:
			return np, true
		default:
			p = np
		}
	}
	return noPlayer, false
}

func (g game) moveThiefFinishTurn(c *gin.Context) (game, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, err := g.validateMoveThiefFinishTurn(c)
	if err != nil {
		return g, err
	}

	g = g.endOfTurnUpdateFor(cp)

	np, found := g.moveThiefNextPlayer(cp)
	if !found {
		// game Over
		g.Phase = phaseGameOver
		return g, nil
	}

	// If game did not end, select next player and continue moving thieves.

	g.CPUserIndices = []int{np.id}
	g.Turn++

	g.Phase = phasePlayCard
	cu := user.Current(c)
	g = g.updateClickablesFor(cu)

	// log.Debugf("g.Animations: %#v", g.Animations)
	// log.Debugf("g.Log: %#v", g.Log)

	g = g.beginningOfTurnReset(np)

	return g, nil
}

func (g game) validateMoveThiefFinishTurn(c *gin.Context) (player, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, err := g.validateFinishTurn(c)
	switch {
	case err != nil:
		return noPlayer, err
	case g.Phase != phaseClaimItem:
		return noPlayer, errWrongPhase
	default:
		return cp, nil
	}
}
