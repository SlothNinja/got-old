package main

import (
	"fmt"

	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// FinishTurn provides a handler for finishing a turn.
func (g game) FinishTurn(c *gin.Context) (game, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	switch g.Phase {
	case phasePlaceThieves:
		return g.placeThievesFinishTurn(c)
	case phaseClaimItem:
		return g.moveThiefFinishTurn(c)
	default:
		return g, errors.WithMessagef(errValidation,
			"incorrect phase %q for finishing turn", phaseName(g.Phase))
	}
}

func (g game) validateFinishTurn(c *gin.Context) (player, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cu, found := user.Current(c)
	if !found {
		return player{}, fmt.Errorf("unable to find current user")
	}

	cp, found := g.currentPlayerFor(cu)
	switch {
	case !found:
		return player{}, errors.WithMessage(errValidation, "current player not found")
	case !g.CPorAdmin(cp.id, cu):
		return player{}, errors.WithMessage(errValidation,
			"only the current player can perform the selected action")
	case !cp.performedAction:
		return player{}, errors.WithMessage(errValidation, "you have yet to perform an action")
	default:
		return cp, nil
	}
}

func (g game) nextPlayer(cp player) (player, bool) {
	index, found := playerFindIndex(cp, g.players)
	if !found {
		return player{}, false
	}
	return playerByIndex(index+1, g.players), true
}

func (g game) previousPlayer(cp player) (player, bool) {
	index, found := playerFindIndex(cp, g.players)
	if !found {
		return player{}, false
	}
	return playerByIndex(index-1, g.players), true
}

func (g game) placeThievesFinishTurn(c *gin.Context) (game, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, err := g.validatePlaceThievesFinishTurn(c)
	if err != nil {
		return g, err
	}

	np, found := g.previousPlayer(cp)
	if !found {
		return g, errors.New("player not found")
	}

	newTurn := cp.Equal(g.players[0])
	if newTurn {
		g.Turn++
	}

	numThieves := 3
	if g.TwoThiefVariant {
		numThieves = 2
	}

	if newTurn && g.Turn > numThieves {
		g = g.startCardPlay()
		np = cp
	}

	g.CPUserIndices = []int{np.id}

	g.Stack = g.Commit()
	g = g.beginningOfTurnReset(np)

	cu, _ := user.Current(c)
	g = g.updateClickablesFor(cu)

	return g, nil
}

func (g game) validatePlaceThievesFinishTurn(c *gin.Context) (player, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, err := g.validateFinishTurn(c)
	switch {
	case err != nil:
		return player{}, err
	case g.Phase != phasePlaceThieves:
		return player{}, errors.WithMessage(errValidation, "wrong phase for selected action")
	default:
		return cp, nil
	}
}

func (g game) moveThiefNextPlayer(p player) (player, bool) {
	for !allPassed(g.players) {
		np, found := g.nextPlayer(p)
		switch {
		case !found:
			return player{}, false
		case !np.passed:
			return np, true
		default:
			p = np
		}
	}
	return player{}, false
}

func (g game) moveThiefFinishTurn(c *gin.Context) (game, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, err := g.validateMoveThiefFinishTurn(c)
	if err != nil {
		return g, err
	}

	g = g.endOfTurnUpdateFor(cp)

	np, found := g.moveThiefNextPlayer(cp)
	if !found {
		// game Over
		g.Phase = phaseGameOver
		return g, nil
	}

	// If game did not end, select next player and continue moving thieves.

	g.CPUserIndices = []int{np.id}
	g.Turn++

	g.Phase = phasePlayCard
	cu, _ := user.Current(c)
	g = g.updateClickablesFor(cu)

	// log.Debugf("g.Animations: %#v", g.Animations)
	// log.Debugf("g.Log: %#v", g.Log)

	g = g.beginningOfTurnReset(np)

	return g, nil
}

func (g game) validateMoveThiefFinishTurn(c *gin.Context) (player, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, err := g.validateFinishTurn(c)
	switch {
	case err != nil:
		return player{}, err
	case g.Phase != phaseClaimItem:
		return player{}, errors.WithMessage(errValidation, "wrong phase for selected action")
	default:
		return cp, nil
	}
}
