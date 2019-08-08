package main

import (
	"net/http"
	"time"

	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/status"
	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func (s server) finish(param, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

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

		g, err = g.FinishTurn(c)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		var rs []result

		if g.Phase == phaseGameOver {
			g.Status = status.Ending
			rs, err = g.endgame(c)
			if err != nil {
				log.Errorf(err.Error())
				c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
				return
			}
		}

		t := time.Now()
		g.UpdatedAt = t

		// h := newHistory(g)
		// e := newHeaderEntity(g.ID())
		// e.Header, e.UpdatedAt = g.Header, g.UpdatedAt

		// ks := []*datastore.Key{g.Key, h.Key, e.Key}
		// es := []interface{}{&g, &h, &e}
		_, err = s.RunInTransaction(c, func(tx *datastore.Transaction) error {
			_, err = tx.PutMulti(g.withHistory())
			if err != nil {
				return err
			}
			err = sendEndGameNotifications(c, g.Header, rs)
			if err != nil {
				log.Warningf(err.Error())
			}
			return nil
		})

		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"game": g})
	}
}
