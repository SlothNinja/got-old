package main

import (
	"net/http"

	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/status"
	"bitbucket.org/SlothNinja/store"
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

func jsonIndexAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		s := status.StatusFromParam(c)
		q := datastore.
			NewQuery("Header").
			Filter("Status=", int(s)).
			Order("-UpdatedAt")

		client, err := store.New(c)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		var es []headerEntity
		ks, err := client.GetAll(c, q, &es)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}
		for i := range ks {
			log.Debugf("ks[%d]: %#v", i, ks[i])
		}
		for i := range es {
			log.Debugf("es[%d]: %#v", i, es[i])
		}

		// hs = make([]Header, len(ks))
		cu, found := user.Current(c)
		if !found {
			log.Errorf("unable to find current user")
			c.JSON(http.StatusOK, gin.H{"message": "unable to find current user"})
			return
		}
		// err = client.GetMulti(c, ks, hs)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"headers": es, "cu": cu})
			return
		}
		c.JSON(http.StatusBadRequest, struct {
			Error string
		}{err.Error()})
	}
}
