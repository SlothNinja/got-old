package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.chromium.org/gae/service/info"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"

	"cloud.google.com/go/datastore"
)

const (
	userNewPath  = "/#/new"
	userShowPath = "/#/show/"
)

// const (
// 	welcomePath = "/welcome"
// 	userNewPath = "/user"
// 	homePath    = "/"
// )
//
// func Index(c *gin.Context) {
// 	c.HTML(http.StatusOK, "user/index", gin.H{
// 		"Context":   c,
// 		"VersionID": info.VersionID(c),
// 		"CUser":     user.CurrentFrom(c),
// 	})
// }
//
// func Show(c *gin.Context) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Entering")
//
// 	u, err := user.ByParam(c, "uid")
// 	if err != nil {
// 		log.Errorf(err.Error())
// 		return
// 	}
//
// 	s, err := stats.ByUser(c, u)
// 	if err != nil {
// 		log.Errorf(err.Error())
// 	}
//
// 	cu := user.CurrentFrom2(c)
// 	c.HTML(http.StatusOK, "user/show", gin.H{
// 		"Context":   c,
// 		"VersionID": "",
// 		"User":      u,
// 		"CUser":     cu,
// 		"IsAdmin":   cu.Admin,
// 		"Stats":     s,
// 	})
// }
//
// func Edit(c *gin.Context) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Entering")
//
// 	u, err := user.ByParam(c, "uid")
// 	if err != nil {
// 		log.Errorf(err.Error())
// 		return
// 	}
//
// 	s, err := stats.ByUser(c, u)
// 	if err != nil {
// 		log.Errorf(err.Error())
// 	}
//
// 	cu := user.CurrentFrom2(c)
// 	c.HTML(http.StatusOK, "user/edit", gin.H{
// 		"Context":   c,
// 		"VersionID": "",
// 		"User":      u,
// 		"CUser":     cu,
// 		"IsAdmin":   cu.Admin,
// 		"Stats":     s,
// 	})
// }

//func Remote(c *restful.Context, render render.Render, params martini.Params) {
//	if u, err := user.ByGoogleID(c, params["uid"]); err == nil {
//		render.JSON(http.StatusOK, u)
//	} else {
//		render.HTML(http.StatusGone, "", "")
//	}
//}

// type jUserIndex struct {
// 	Data            []*jUser `json:"data"`
// 	Draw            int      `json:"draw"`
// 	RecordsTotal    int      `json:"recordsTotal"`
// 	RecordsFiltered int      `json:"recordsFiltered"`
// }
//
// type omit *struct{}
//
// type jUser struct {
// 	IntID         int64         `json:"id"`
// 	StringID      string        `json:"sid"`
// 	OldID         int64         `json:"oldid"`
// 	GoogleID      string        `json:"googleid"`
// 	Name          string        `json:"name"`
// 	Email         string        `json:"email"`
// 	Gravatar      template.HTML `json:"gravatar"`
// 	JoinedAt      time.Time     `json:"joinedAt"`
// 	UpdatedAt     time.Time     `json:"updatedAt"`
// 	OmitCreatedAt omit          `json:"createdat,omitempty"`
// 	OmitUpdatedAt omit          `json:"updatedat,omitempty"`
// }
//
// func toUserTable(c *gin.Context, us []user.User2) (table *jUserIndex, err error) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	table = new(jUserIndex)
// 	l := len(us)
// 	table.Data = make([]*jUser, l)
//
// 	for i, u := range us {
// 		table.Data[i] = &jUser{
// 			StringID:  u.ID(),
// 			OldID:     0,
// 			GoogleID:  u.GoogleID,
// 			Name:      u.Name,
// 			Email:     u.Email,
// 			Gravatar:  user.Gravatar2(u),
// 			JoinedAt:  u.CreatedAt,
// 			UpdatedAt: u.UpdatedAt,
// 		}
// 	}
//
// 	if draw, err := strconv.Atoi(c.PostForm("draw")); err != nil {
// 		return nil, err
// 	} else {
// 		table.Draw = draw
// 	}
// 	table.RecordsTotal = user.CountFrom(c)
// 	table.RecordsFiltered = user.CountFrom(c)
// 	return
// }
//
// func JSON(c *gin.Context) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	us := user.UsersFrom(c)
// 	if data, err := toUserTable(c, us); err != nil {
// 		c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
// 	} else {
// 		c.JSON(http.StatusOK, data)
// 	}
// }

// func NewAction(c *gin.Context) {
// 	u := user.New0(c)
// 	gu := user.GUserFrom(c)
// 	if gu == nil {
// 		restful.AddErrorf(c, "You must be logged in to access this page.")
// 		c.Redirect(http.StatusSeeOther, welcomePath)
// 		return
// 	}
//
// 	u.Name = strings.Split(gu.Email, "@")[0]
// 	u.LCName = strings.ToLower(u.Name)
// 	u.Email = gu.Email
//
// 	c.HTML(http.StatusOK, "user/new", gin.H{
// 		"Context": c,
// 		"User":    user.FromGUser(c, user.GUserFrom(c)),
// 	})
// }

func newUserAction(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu, found := user.Current(c)
		if !found {
			log.Errorf("unable to find current user")
			c.JSON(http.StatusOK, gin.H{"dev": isDev(), "message": "unable to find current user"})
			return
		}

		if !cu.Admin {
			jsonMsg(c, "You already have an account.")
			return
		}

		s := sessions.Default(c)
		token, ok := user.SessionTokenFrom(s)
		if !ok {
			log.Errorf("Missing SessionToken")
			jsonMsg(c, "Unexpected error. Try again.")
			return
		}

		log.Debugf("token: %#v", token)
		u := user.New2(token.ID)
		u.Name = user.Name(token.Email)
		u.Email = token.Email
		u.EmailReminders = true
		u.EmailNotifications = true
		u.GravType = user.GTMonster
		c.JSON(http.StatusOK, gin.H{"u": u})
	}
}

func current(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu, found := user.Current(c)
		if !found {
			c.JSON(http.StatusOK, gin.H{"dev": isDev(), "message": "unable to find current user"})
		}

		c.JSON(http.StatusOK, gin.H{"dev": isDev(), "cu": cu})
	}
}

func jsonUser(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		id := c.Param("id")
		u, err := user.ByID(c, id)
		if err != nil {
			log.Warningf("ByID err: %s", err)
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}

		c.JSON(http.StatusOK, gin.H{"u": u})
	}
}

// func NewDevAction(c *gin.Context) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	s := sessions.Default(c)
// 	token, ok := user.SessionTokenFrom(s)
// 	if !ok {
// 		log.Errorf("Missing SessionToken")
// 		c.Redirect(http.StatusSeeOther, homePath)
// 	}
//
// 	log.Debugf("token: %#v", token)
// 	u := user.New2(token.ID)
// 	u.Name = user.Name(token.Email)
// 	c.HTML(http.StatusOK, "user/new_dev", gin.H{
// 		"Context": c,
// 		"User":    u,
// 	})
// }
//
func createUser(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu, found := user.Current(c)
		if !found {
			c.JSON(http.StatusOK, gin.H{"message": "unable to find current user"})
			return
		}

		if !cu.Admin {
			c.JSON(http.StatusOK, gin.H{"message": "must be logged-in to create account"})
			return
		}

		session := sessions.Default(c)
		token, ok := user.SessionTokenFrom(session)
		if !ok {
			log.Errorf("Missing SessionToken")
			jsonMsg(c, "Must be logged-in to create account.")
			return
		}

		u, err := user.ByID(c, token.ID)
		if err == nil {
			jsonMsg(c, "Account already exists for id: %s", token.ID)
			return
		}

		if err != datastore.ErrNoSuchEntity {
			jsonMsg(c, "Unexpected error. Try again.")
			return
		}

		u, err = fromJSON(c, token.ID, token.Email)
		if err != nil {
			log.Errorf(err.Error())
			jsonMsg(c, "Unexpected error. Try again.")
			return
		}

		id0, id1, gid, joinedAt, found, err := user.ByOldEntries(c, u.Email)
		if err != nil {
			log.Errorf("user.ByOldEntries error: %s", err)
			jsonMsg(c, "Unexpected error. Try again.")
			return
		}

		log.Debugf("found: %#v", found)
		if found {
			u.User0ID = id0
			u.User1ID = id1
			u.GoogleID = gid
			u.JoinedAt = joinedAt
		}

		client, err := user.Client(c)
		if err != nil {
			log.Errorf(err.Error())
			jsonMsg(c, "Unexpected error. Try again.")
			return
		}

		t := time.Now()
		u.UpdatedAt, u.CreatedAt = t, t
		if u.JoinedAt.IsZero() {
			u.JoinedAt = t
		}
		_, err = client.Put(c, u.Key, u)
		if err != nil {
			log.Errorf("datastore.Put error: %v", err)
			jsonMsg(c, "Unexpected error. Try again.")
			return
		}

		log.Debugf("put user: %#v", u)

		err = u.To(session)
		if err != nil {
			log.Errorf("session.Save error: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{"u": u})
	}
}

func fromJSON(c *gin.Context, id, email string) (user.User2, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	jData := struct {
		Name               string `json:"name"`
		EmailNotifications bool   `json:"emailnotifications"`
		GravType           string `json:"gravtype"`
	}{}

	err := c.ShouldBindJSON(&jData)
	if err != nil {
		return user.User2{}, err
	}

	u := user.New2(id)
	u.Name = strings.TrimSpace(jData.Name)
	u.LCName = strings.ToLower(jData.Name)
	u.Email = email
	u.EmailNotifications = jData.EmailNotifications
	u.EmailReminders = true
	u.GravType = jData.GravType
	return u, nil
}

// func Create(prefix string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		u := user.FromGUser(c, user.GUserFrom(c))
// 		switch existing, err := user.ByGoogleID(c, u.GoogleID); {
// 		case err == user.ErrNotFound:
// 		case err != nil:
// 			restful.AddErrorf(c, err.Error())
// 			c.Redirect(http.StatusSeeOther, userNewPath)
// 			return
// 		case existing != nil:
// 			restful.AddErrorf(c, "You already have an account.")
// 			c.Redirect(http.StatusSeeOther, homePath)
// 			return
// 		default:
// 			log.Error("Unexpected result for user.Create. err: %v existing: %v", err, existing)
// 			c.Redirect(http.StatusSeeOther, userNewPath)
// 			return
// 		}
//
// 		// Fell through 'switch' thus err == user.ErrNotFound
// 		u.Name = strings.Split(c.PostForm("user-name"), "@")[0]
// 		u.LCName = strings.ToLower(u.Name)
// 		//u.Key = user.NewKey(c, 0)
//
// 		n := name.New()
// 		if !name.IsUnique(c, u.LCName) {
// 			restful.AddErrorf(c, "%q is not a unique user name.", u.LCName)
// 			c.Redirect(http.StatusSeeOther, userNewPath)
// 			return
// 		}
//
// 		n.GoogleID = u.GoogleID
// 		n.ID = u.LCName
//
// 		err := datastore.RunInTransaction(c, func(tc *gin.Context) (terr error) {
// 			entities := []interface{}{u, n}
// 			if terr = datastore.Put(tc, entities); terr != nil {
// 				return
// 			}
// 			nu := user.ToNUser(c, u)
// 			return datastore.Put(tc, nu)
// 		}, &datastore.TransactionOptions{XG: true})
//
// 		if err != nil {
// 			log.Error("User/Controller#Create datastore.RunInTransaction Error: %v", err)
// 			c.Redirect(http.StatusSeeOther, homePath)
// 			return
// 		}
//
// 		c.Redirect(http.StatusSeeOther, showPath(prefix, u.Key.Name))
// 	}
// }

func showUserPath(prefix string, uid string) string {
	return fmt.Sprint(prefix, userShowPath, uid)
}

//func SendTestMessage(c *restful.Context, render render.Render, routes martini.Routes, params martini.Params) {
//	u := user.Fetched(c)
//	m := new(xmpp.Message)
//	m.To = []string{u.Email}
//	m.Body = fmt.Sprintf("Test message from SlothNinja Games for %s", u.Name)
//	send.XMPP(c, m)
//	c.AddNoticef("Test IM sent to %s", u.Name)
//	render.Redirect(routes.URLFor("user_show", params["uid"]), http.StatusSeeOther)
//}
//
//func SendIMInvite(c *restful.Context, render render.Render, routes martini.Routes, params martini.Params) {
//	u := user.Fetched(c)
//	send.Invite(c, u.Email)
//	c.AddNoticef("IM Invite sent to %s", u.Name)
//	render.Redirect(routes.URLFor("user_show", params["uid"]), http.StatusSeeOther)
//}
func update(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		cu, found := user.Current(c)
		if !found {
			jsonMsg(c, "Must be logged-in to edit account.")
			return
		}

		uid := c.Param("uid")
		if uid != cu.ID() && !cu.Admin {
			jsonMsg(c, "You can only edit your own account.")
			return
		}

		u, err := user.ByID(c, uid)
		if err != nil {
			log.Errorf(err.Error())
			jsonMsg(c, "Unexpected error. Try again.")
			return
		}

		u2, err := fromJSON(c, u.ID(), u.Email)
		if err != nil {
			log.Errorf(err.Error())
			jsonMsg(c, "Unexpected error. Try again.")
			return
		}

		uniq, err := uniqueName(c, u2)
		if err != nil {
			log.Errorf(err.Error())
			jsonMsg(c, "Unexpected error. Try again.")
			return
		}

		if !uniq {
			jsonMsg(c, "Screen Name: %s already in use.", u2.Name)
			return
		}

		// oldName := name.New()
		// oldName.ID = u.LCName
		// if err := u.Update(c, u); err != nil {
		// 	log.Errorf("User/Controller#Update u.update Error: %s", err)
		// 	route := fmt.Sprintf("/user/show/%s", c.Param("uid"))
		// 	c.Redirect(http.StatusSeeOther, route)
		// 	return
		// }
		// newName := name.New()
		// newName.GoogleID = u.GoogleID
		// newName.ID = u.LCName

		// log.Debug("Before datastore.RunInTransaction")
		// err = datastore.RunInTransaction(c, func(tc *gin.Context) (err error) {
		// 	nu := user.ToNUser(c, u)
		// 	entities := []interface{}{u, nu, newName, oldName}
		// 	if err = datastore.Put(tc, entities); err != nil {
		// 		return
		// 	}

		// 	return datastore.Delete(tc, oldName)
		// }, &datastore.TransactionOptions{XG: true})

		client, err := user.Client(c)
		if err != nil {
			log.Errorf(err.Error())
			return
		}

		_, err = client.Put(c, u2.Key, u2)
		if err != nil {
			log.Errorf(err.Error())
			// route := fmt.Sprintf("/user/show/%s", c.Param("uid"))
			// c.Redirect(http.StatusSeeOther, route)
			return
		}

		c.JSON(http.StatusOK, gin.H{"cu": u2})
	}
}

func jsonMsg(c *gin.Context, format string, args ...interface{}) {
	c.JSON(http.StatusOK, gin.H{"dev": isDev(), "msg": fmt.Sprintf(format, args...)})
}

func uniqueName(c *gin.Context, u1 user.User2) (bool, error) {
	u2, err := user.ByLCName(c, u1.LCName)
	if err != nil {
		return false, err
	}

	if u1.LCName == u2.LCName && u1.ID() != u2.ID() {
		return false, nil
	}
	return true, nil
}

// func GamesIndex(c *gin.Context) {
// 	log.Debug("Entering")
// 	defer log.Debug("Exiting")
//
// 	if status := game.StatusFrom(c); status != game.NoStatus {
// 		c.HTML(200, "shared/games_index", gin.H{})
// 	} else {
// 		c.HTML(200, "user/games_index", gin.H{})
// 	}
// }

func Login(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	session := sessions.Default(c)
	state := randToken()
	session.Set("state", state)
	session.Save()

	if v := os.Getenv("DEV_LOGIN"); v == "true" {
		c.HTML(http.StatusOK, "user/login", gin.H{
			"Context":   c,
			"VersionID": info.VersionID(c),
		})
	} else {
		c.Redirect(http.StatusSeeOther, getLoginURL(c, state))
	}
}

// func sessionLogin(session sessions.Session) bool {
// 	uinfo, ok := user.InfoFrom(session)
// 	if !ok {
// 		return false
// 	}
//
// 	if uinfo.LoggedIn {
// 		return true
// 	}
//
// 	uinfo.LoggedIn = true
// 	err := uinfo.To(session)
// 	return err == nil
// }

func Logout(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.Redirect(http.StatusSeeOther, "/")
}

func getLoginURL(c *gin.Context, state string) string {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	return oauth2Config(c, scopes()...).AuthCodeURL(state)
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func DevAuth(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	email := c.PostForm("auth-email")
	log.Debugf("email: %v", email)

	// Handle the exchange code to initiate a transport.
	//	retrievedState := session.Get("state")
	//	if retrievedState != c.Query("state") {
	//		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
	//		return
	//	}
	//
	//	log.Debug("retrievedState: %#v", retrievedState)
	//	ac := appengine.NewContext(c.Request)
	//	tok, err := conf.Exchange(ac, c.Query("code"))
	//	if err != nil {
	//		log.Error("tok error: %#v\n %s", err, err)
	//		c.AbortWithError(http.StatusBadRequest, err)
	//		return
	//	}
	//
	//	log.Debug("tok: %#v", tok)
	//
	//	client := conf.Client(ac, tok)
	//	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	//	if err != nil {
	//		c.AbortWithError(http.StatusBadRequest, err)
	//		return
	//	}
	//	defer resp.Body.Close()
	//	body, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		c.AbortWithError(http.StatusBadRequest, err)
	//		return
	//	}
	//
	//	info := user.Info{}
	//	var b binding.BindingBody = binding.JSON
	//	err = b.BindBody(body, &info)
	//	if err != nil {
	//		log.Error("BindBody error: %v", err)
	//		c.AbortWithError(http.StatusBadRequest, err)
	//		return
	//	}
	//	log.Debug("info: %#v", info)
	//
	u := user.New2(user.ID(email))
	u.LCName = user.LCName(email)
	session := sessions.Default(c)
	err := u.To(session)
	if err != nil {
		log.Errorf("session.Save error: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Debugf("session saved")

	if _, err = user.ByID(c, u.ID()); err != nil {
		if err != user.ErrNotFound {
			log.Errorf("err: %#v user.ErrNotFound: %#v", err, user.ErrNotFound)
			log.Errorf("err T: %T user.ErrNotFound: %T", err, user.ErrNotFound)
			log.Errorf("user.ByGoogleID error: %v", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.Redirect(http.StatusSeeOther, "/user/new_dev")
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func Auth(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}

	log.Debugf("retrievedState: %#v", retrievedState)
	ac := appengine.NewContext(c.Request)
	conf := oauth2Config(c, scopes()...)
	tok, err := conf.Exchange(ac, c.Query("code"))
	if err != nil {
		log.Errorf("tok error: %#v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	log.Debugf("tok: %#v", tok)

	client := conf.Client(ac, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Debugf("body: %s", body)

	uinfo := user.Info{}
	var b binding.BindingBody = binding.JSON
	err = b.BindBody(body, &uinfo)
	if err != nil {
		log.Errorf("BindBody error: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Debugf("info: %#v", uinfo)

	id := user.ID(uinfo.Sub)
	u, err := user.ByID(c, id)
	if err == datastore.ErrNoSuchEntity {
		u = user.New2(id)
		u.Email = uinfo.Email
		err = u.To(session)
		if err != nil {
			log.Errorf("session.Save error: %v", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		log.Debugf("session saved")
		c.Redirect(http.StatusSeeOther, userNewPath)
		return
	}

	if err != nil {
		log.Errorf("user.ByID => \n\t id: %s\n\t error: %s", id, err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = u.To(session)
	if err != nil {
		log.Errorf("session.Save error: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Debugf("session saved")

	c.Redirect(http.StatusSeeOther, "/")
}

func scopes() []string {
	return []string{"email", "profile", "openid"}
}

func oauth2Config(c *gin.Context, scopes ...string) *oauth2.Config {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	conf := &oauth2.Config{
		ClientID:     "435340145701-t5o50sjq7hsbilopgreobhvrv30e1tj4.apps.googleusercontent.com",
		ClientSecret: "Fe5f-Ht1V5_GohDEOS_TQOVc",
		Endpoint:     google.Endpoint,
		Scopes:       scopes,
		RedirectURL:  "https://user.slothninja.com/auth",
	}
	if isDev() {
		port := getPort()
		conf.RedirectURL = fmt.Sprintf("http://user.slothninja.com:%s/auth", port)
	}
	return conf
}
