package main

import (
	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const playCardID = "play-card"

func (g game) startCardPlay() game {
	g.Phase = phasePlayCard

	np := playerByIndex(0, g.players)
	g = g.beginningOfTurnReset(np)

	g.CPUserIndices = []int{np.ID}
	return g
}

func (g game) PlayCard(c *gin.Context) (game, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, index, err := g.validatePlayCard(c)
	if err != nil {
		return g, err
	}

	var cd card
	cp.Hand, cd = drawFrom(index, cp.Hand)
	cp.DiscardPile = append(cp.DiscardPile, cd)
	g.updatePlayer(cp)

	if cd.kind == cdJewels {
		g.playedCard = g.jewels
	} else {
		g.playedCard = cd
	}

	g.Log = nil
	g.Log = append(g.Log, logEntry{
		"template": playCardID,
		"pid":      cp.ID,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"card":     cd,
	})

	g.Stack = g.Stack.Update()

	g.Phase = phaseSelectThief
	cu, _ := user.Current(c)
	g = g.updateClickablesFor(cu)

	return g, nil
}

func (g game) validatePlayCard(c *gin.Context) (player, int, error) {
	cp, err := g.validatePlayerAction(c)
	switch {
	case err != nil:
		return player{}, 0, err
	case g.Phase != phasePlayCard:
		return player{}, 0, errors.Wrap(errValidation, "wrong phase for playing a card")
	default:
		index, err := getIndex(c, cp.Hand)
		if err != nil {
			return player{}, 0, err
		}
		return cp, index, nil
	}
}
