package main

import (
	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const placeThiefID = "place-thief"

func (g game) PlaceThief(c *gin.Context) (game, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, a, err := g.validatePlaceThief(c)
	if err != nil {
		return g, err
	}

	g.Log = nil
	g.Log = append(g.Log, logEntry{
		"template": placeThiefID,
		"pid":      cp.ID,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"area":     a,
	})

	// Update game state to reflect placed thief.
	cp.PerformedAction = true
	cp.Score += a.card.value()
	g = g.updatePlayer(cp)

	g, a = g.placeThief(cp, a)

	g.Stack = g.Stack.Update()
	cu, _ := user.Current(c)
	g.updateClickablesFor(cu)

	return g, nil
}

func (g game) placeThief(p player, a area) (game, area) {
	a.thief.pid = p.ID
	return g.updateArea(a), a
}

func (g game) placeThieves() game {
	g.Phase = phasePlaceThieves
	return g
}

func (g game) validatePlaceThief(c *gin.Context) (player, area, error) {
	cp, err := g.validatePlayerAction(c)
	if err != nil {
		return player{}, area{}, err
	}

	a, err := g.getArea(c)
	switch {
	case err != nil:
		return player{}, area{}, err
	case a.card.kind == cdNone:
		return player{}, area{}, errors.Wrap(errValidation, "selected area has no card to claim")
	case a.thief.pid != pidNone:
		return player{}, area{}, errors.Wrap(errValidation, "selected area already has a thief")
	case g.Phase != phasePlaceThieves:
		return player{}, area{}, errors.Wrap(errValidation, "wrong phase for placing thieves")
	default:
		return cp, a, nil
	}
}
