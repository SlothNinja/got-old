package main

import (
	"fmt"

	"github.com/SlothNinja/log"
	"github.com/SlothNinja/user"

	"github.com/gin-gonic/gin"
)

const moveThiefID = "move-thief"
const bumpedThiefID = "bumped-thief"

func (g game) startMoveThief() game {
	g.Phase = phaseMoveThief
	return g
}

func (s server) moveThief() gin.HandlerFunc {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return s.update(param, (game).moveThiefAction)
}

func (g game) moveThiefAction(c *gin.Context) (game, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, from, to, cd, err := g.validateMoveThief(c)
	if err != nil {
		return g, err
	}

	g.Log = append(g.Log, logEntry{
		"template": moveThiefID,
		"pid":      cp.id,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"card":     cd,
		"from":     from,
		"to":       to,
	})

	cp.score += to.card.value()

	switch cd.kind {
	case cdSword:
		g, cp, from, to = g.swordMove(cp, from, to)
	case cdTurban:
		g, cp, from, to = g.turbanMove(cp, from, to)
	case cdCoins:
		g, cp, from, to = g.coinMove(cp, from, to)
	default:
		g, cp, from, to = g.defaultMove(cp, from, to)
	}

	cu := user.Current(c)
	g = g.updateClickablesFor(cu)
	g.Stack = g.Stack.Update()
	return g, nil
}

func (g game) swordMove(cp player, from, to area) (game, player, area, area) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	p2 := playerByID(to.thief.pid, g.players)
	bumpedTo := g.grid.bumpedTo(from, to)
	g.grid, to, bumpedTo = g.grid.moveThief(to, bumpedTo)
	p2.score += bumpedTo.card.value() - to.card.value()
	g = g.updatePlayer(p2)

	// Move thief
	g.grid, from, to = g.grid.moveThief(from, to)

	// Claim Item
	g, cp, from = g.claimItemFor(cp, from)
	if !g.toHand() {
		cp, _, _ = cp.draw()
	}

	cp.performedAction = true
	return g.updatePlayer(cp), cp, from, to
}

// func (g grid) swordMove(cp player, from, to area) grid {
// 	log.Debugf(msgEnter)
// 	defer log.Debugf(msgExit)
//
// 	bumpedTo := g.bumpedTo(from, to)
// 	g, _, _ = g.moveThief(to, bumpedTo)
//
// 	// Move thief
// 	g, _, _ = g.moveThief(from, to)
//
// 	return g
// }

func (g game) turbanMove(cp player, from, to area) (game, player, area, area) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	if g.stepped == 0 {
		g.stepped = 1
		g.selectedAreaID = to.areaID
		g = g.startMoveThief()

		g, cp, from, to = g.defaultMove(cp, from, to)

		// Revised defaultMove
		cp.performedAction = false
		g.Phase = phaseMoveThief
		return g, cp, from, to
	}

	g.stepped = 2
	return g.defaultMove(cp, from, to)
}

func (g game) coinMove(cp player, from, to area) (game, player, area, area) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	g, cp, from, to = g.defaultMove(cp, from, to)
	cp, _, _ = cp.draw()
	return g, cp, from, to
}

func (g grid) removeThiefFrom(a area) (grid, area) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	a.thief.pid = noPID
	return g.updateArea(a), a
}

func (g grid) moveThief(from, to area) (grid, area, area) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	pid := from.thief.pid
	g, from = g.removeThiefFrom(from)
	to.thief.pid, to.thief.from = pid, from.areaID
	return g.updateArea(to), from, to
}

func (g game) defaultMove(cp player, from, to area) (game, player, area, area) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	// Move thief
	g.grid, from, to = g.grid.moveThief(from, to)

	// Claim Item
	g, cp, from = g.claimItemFor(cp, from)
	if !g.toHand() {
		cp, _, _ = cp.draw()
	}

	cp.performedAction = true
	return g.updatePlayer(cp), cp, from, to
}

func (g game) validateMoveThief(c *gin.Context) (player, area, area, card, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cp, err := g.validatePlayerAction(c)
	if err != nil {
		return noPlayer, noArea, noArea, noCard, err
	}

	to, err := g.getArea(c)
	if err != nil {
		return noPlayer, noArea, noArea, noCard, err
	}

	from := g.grid.area(g.selectedAreaID.row, g.selectedAreaID.column)
	if from == noArea {
		return noPlayer, noArea, noArea, noCard, fmt.Errorf("selected thief area not found: %w", errValidation)
	}

	cd := g.playedCard
	log.Debugf("g.playedCard: %#v", g.playedCard)
	switch {
	case from.thief.pid != cp.id:
		return noPlayer, noArea, noArea, noCard,
			fmt.Errorf("selected thief of another player: %w", errValidation)
	case cd.kind == cdNone:
		return noPlayer, noArea, noArea, noCard,
			fmt.Errorf("you must play card before moving thief: %w", errValidation)
	case cd.kind == cdGuard:
		return noPlayer, noArea, noArea, noCard,
			fmt.Errorf("played card does not permit moving selected thief to selected area: %w", errValidation)
	case (cd.kind == cdLamp || cd.kind == cdSLamp) && g.grid.isLampMove(from, to):
		return cp, from, to, cd, nil
	case (cd.kind == cdCamel || cd.kind == cdSCamel) && g.grid.isCamelMove(from, to):
		return cp, from, to, cd, nil
	case cd.kind == cdCoins && g.isCoinsArea(to):
		return cp, from, to, cd, nil
	case cd.kind == cdSword && g.grid.isSwordMoveFor(cp, from, to):
		return cp, from, to, cd, nil
	case cd.kind == cdCarpet && g.grid.isCarpetMove(from, to):
		return cp, from, to, cd, nil
	case cd.kind == cdTurban && g.stepped == 0 && g.grid.isTurban0Move(from, to):
		return cp, from, to, cd, nil
	case cd.kind == cdTurban && g.stepped == 1 && g.isTurban1Area(to):
		return cp, from, to, cd, nil
	default:
		return noPlayer, noArea, noArea, noCard,
			fmt.Errorf("played card does not permit moving selected thief to selected area: %w", errValidation)
	}
}

// bumpedTo assumes bumping a thief by moving thief 'from' area 'to' area is valid.
// calling validateMoveThief before bumpedTo ensures bumping a thief by moving thief 'from' area 'to' area is valid.
func (g grid) bumpedTo(from, to area) area {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	switch {
	case from.row > to.row:
		return g.area(to.row-1, from.column)
	case from.row < to.row:
		return g.area(to.row+1, from.column)
	case from.column > to.column:
		return g.area(from.row, to.column-1)
	case from.column < to.column:
		return g.area(from.row, to.column+1)
	}
	return noArea
}
