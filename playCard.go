package main

import (
	"fmt"

	"github.com/SlothNinja/log"
	"github.com/SlothNinja/user"

	"github.com/gin-gonic/gin"
)

const playCardID = "play-card"

func (g game) startCardPlay() game {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	g.Phase = phasePlayCard

	np := playerByIndex(0, g.players)
	g = g.beginningOfTurnReset(np)

	g.CPUserIndices = []int{np.id}
	return g
}

func (g game) PlayCard(c *gin.Context) (game, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, index, err := g.validatePlayCard(c)
	if err != nil {
		return g, err
	}

	var cd card
	cp.hand, cd = drawFrom(index, cp.hand)
	cp.discardPile = append(cp.discardPile, cd)
	g.updatePlayer(cp)

	if cd.kind == cdJewels {
		g.playedCard = g.jewels
	} else {
		g.playedCard = cd
	}

	g.Log = nil
	g.Log = append(g.Log, logEntry{
		"template": playCardID,
		"pid":      cp.id,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"card":     cd,
	})

	g.Stack = g.Stack.Update()

	g.Phase = phaseSelectThief
	cu := user.Current(c)
	g = g.updateClickablesFor(cu)

	return g, nil
}

func (g game) validatePlayCard(c *gin.Context) (player, int, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, err := g.validatePlayerAction(c)
	switch {
	case err != nil:
		return noPlayer, 0, err
	case g.Phase != phasePlayCard:
		return noPlayer, 0, fmt.Errorf("wrong phase for playing a card: %w", errValidation)
	default:
		index, err := getIndex(c, cp.hand)
		if err != nil {
			return noPlayer, 0, err
		}
		return cp, index, nil
	}
}
