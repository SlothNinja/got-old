package main

import (
	"github.com/SlothNinja/log"
	"github.com/gin-gonic/gin"
)

const claimItemID = "claim-item"

func (g game) toHand() bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	numThieves := 3
	if g.TwoThiefVariant {
		numThieves = 2
	}

	return g.Turn <= (numThieves+1)*g.NumPlayers
}

func (g grid) removeCardFrom(a area) (grid, area) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	a.card = card{kind: cdNone, facing: cdFaceDown}
	return g.updateArea(a), a
}

func (g game) claimItemFor(cp player, from area) (game, player, area) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	toHand := g.toHand()

	g.Phase = phaseClaimItem
	cd := from.card
	g.grid, from = g.grid.removeCardFrom(from)

	// If first claimed card, place in hand instead of discard pile
	if toHand {
		cd.turn(cdFaceUp)
		cp.hand = append(cp.hand, cd)
	} else {
		cp.discardPile = append(cp.discardPile, cd)
	}

	g.Log = append(g.Log, logEntry{
		"template": claimItemID,
		"pid":      cp.id,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"card":     cd,
		"from":     from,
		"toHand":   toHand,
	})

	return g.updatePlayer(cp), cp, from
}

func (g game) finalClaim(c *gin.Context) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	g.Phase = phaseFinalClaim
	for row := rowA; row <= lastRowFor(g.NumPlayers); row++ {
		for col := col1; col <= col8; col++ {
			a := g.grid.area(row, col)
			if a != noArea {
				p := playerByID(a.thief.pid, g.players)
				if p.id != noPID {
					cd := a.card
					a.card = newCard(cdNone, cdFaceDown)
					a.thief.pid = noPID
					p.discardPile = append([]card{cd}, p.discardPile...)
				}
			}
		}
	}
	for _, p := range g.players {
		p.collectCards()
	}
}
