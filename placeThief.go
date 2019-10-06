package main

import (
	"fmt"

	"github.com/SlothNinja/log"
	"github.com/SlothNinja/user"

	"github.com/gin-gonic/gin"
)

const placeThiefID = "place-thief"

func (s server) placeThief() gin.HandlerFunc {
	return s.update(param, (game).placeThief)
}

func (g game) placeThief(c *gin.Context) (game, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, a, err := g.validatePlaceThief(c)
	if err != nil {
		return g, err
	}

	g.Log = nil
	g.Log = append(g.Log, logEntry{
		"template": placeThiefID,
		"pid":      cp.id,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"area":     a,
	})

	// Update game state to reflect placed thief.
	cp.performedAction = true
	cp.score += a.card.value()
	g = g.updatePlayer(cp)

	g.grid, a = g.grid.placeThiefIn(cp, a)

	g.Stack = g.Stack.Update()
	cu := user.Current(c)
	g.updateClickablesFor(cu)

	return g, nil
}

func (g grid) placeThiefIn(p player, a area) (grid, area) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	a.thief.pid = p.id
	return g.updateArea(a), a
}

func (g game) validatePlaceThief(c *gin.Context) (player, area, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, err := g.validatePlayerAction(c)
	if err != nil {
		return noPlayer, noArea, err
	}

	a, err := g.getArea(c)
	switch {
	case err != nil:
		return noPlayer, noArea, err
	case a.card.kind == cdNone:
		return noPlayer, noArea, fmt.Errorf("selected area has no card to claim: %w", errValidation)
	case a.thief.pid != noPID:
		return noPlayer, noArea, fmt.Errorf("selected area already has a thief: %w", errValidation)
	case g.Phase != phasePlaceThieves:
		return noPlayer, noArea, fmt.Errorf("wrong phase for placing thieves: %w", errValidation)
	default:
		return cp, a, nil
	}
}
