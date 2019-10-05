package main

import (
	"net/http"

	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/status"
	"bitbucket.org/SlothNinja/user"
	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

// func Index() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		log.Debugf("Entering")
// 		defer log.Debugf("Exiting")
//
// 		cu := user.Current(c)
//
// 		switch s := status.StatusFromParam(c); s {
// 		case status.Recruiting:
// 			c.HTML(http.StatusOK, "shared/invitation_index", gin.H{
// 				"Context":   c,
// 				"VersionID": info.VersionID(c),
// 				"CUser":     cu,
// 				"Type":      gtype.GOT.String(),
// 			})
// 		default:
// 			c.HTML(http.StatusOK, "shared/games_index", gin.H{
// 				"Context":   c,
// 				"VersionID": info.VersionID(c),
// 				"CUser":     cu,
// 				"Type":      gtype.GOT.String(),
// 				"Status":    s,
// 			})
// 		}
// 	}
// }

func (s server) jsonIndexAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		stat := status.StatusFromParam(c)
		q := datastore.
			NewQuery("Header").
			Filter("Status=", int(stat)).
			Order("-UpdatedAt")

		s, err := s.init(c)
		if err != nil {
			jerr(c, err)
			return
		}

		var es []headerEntity
		_, err = s.GetAll(c, q, &es)
		if err != nil {
			jerr(c, err)
			return
		}

		cu := user.Current(c)
		if cu == user.None {
			jerr(c, errUserNotFound)
			return
		}

		c.JSON(http.StatusOK, gin.H{"headers": es, "cu": cu})
		return
	}
}
