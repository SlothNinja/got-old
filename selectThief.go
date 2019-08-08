package main

import (
	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const selectThiefID = "select-thief"

func (g game) startSelectThief(c *gin.Context) (tmpl string, err error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	g.Phase = phaseSelectThief
	tmpl = "played_card_update"
	return
}

func (g game) selectThiefIn(a area) game {
	g.selectedAreaID = a.areaID
	return g
}

func (g game) SelectThief(c *gin.Context) (game, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, a, err := g.validateSelectThief(c)
	if err != nil {
		return g, err
	}

	g = g.selectThiefIn(a)
	g = g.startMoveThief()

	g.Log = append(g.Log, logEntry{
		"template": selectThiefID,
		"pid":      cp.ID,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"area":     a,
	})

	cu, _ := user.Current(c)
	g = g.updateClickablesFor(cu)

	g.Stack = g.Stack.Update()
	return g, nil
}

func (g game) validateSelectThief(c *gin.Context) (player, area, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, err := g.validatePlayerAction(c)
	if err != nil {
		return player{}, area{}, err
	}

	a, err := g.getArea(c)
	switch {
	case err != nil:
		return player{}, area{}, err
	case (a.thief.pid != cp.ID):
		return player{}, area{}, errors.WithMessage(errValidation, "selected area does not have one of your thieves")
	default:
		return cp, a, nil
	}
}
