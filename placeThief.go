package main

import (
	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

	g, a = g.placeThiefIn(cp, a)

	g.Stack = g.Stack.Update()
	cu, _ := user.Current(c)
	g.updateClickablesFor(cu)

	return g, nil
}

func (g game) placeThiefIn(p player, a area) (game, area) {
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
		return player{}, area{}, err
	}

	a, err := g.getArea(c)
	switch {
	case err != nil:
		return player{}, area{}, err
	case a.card.kind == cdNone:
		return player{}, area{}, errors.WithMessage(errValidation, "selected area has no card to claim")
	case a.thief.pid != pidNone:
		return player{}, area{}, errors.WithMessage(errValidation, "selected area already has a thief")
	case g.Phase != phasePlaceThieves:
		return player{}, area{}, errors.WithMessage(errValidation, "wrong phase for placing thieves")
	default:
		return cp, a, nil
	}
}
