package main

import (
	"bitbucket.org/SlothNinja/log"
	"github.com/gin-gonic/gin"
)

// func init() {
// 	move.Register(claimItemID, new(claimItemMoveData))
// }

const claimItemID = "claim-item"

func (g game) toHand() bool {
	numThieves := 3
	if g.TwoThiefVariant {
		numThieves = 2
	}

	return g.Turn <= (numThieves+1)*g.NumPlayers
}

func (g game) removeCardFrom(a area) (game, area) {
	a.card = card{kind: cdNone, facing: cdFaceDown}
	return g.updateArea(a), a
}

func (g game) claimItem(a area, cp player) game {
	g.Phase = phaseClaimItem
	cd := a.card
	g, a = g.removeCardFrom(a)

	toHand := g.toHand()
	// If first claimed card, place in hand instead of discard pile
	if toHand {
		cd.turn(cdFaceUp)
		cp.Hand = append(cp.Hand, cd)
	} else {
		cp.DiscardPile = append(cp.DiscardPile, cd)
	}

	g.Log = append(g.Log, logEntry{
		"template": claimItemID,
		"pid":      cp.ID,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"card":     cd,
		"from":     a,
		"toHand":   toHand,
	})

	g = g.updatePlayer(cp)

	log.Debugf("phase: %s", phaseName(g.Phase))

	return g
}

// type claimItemMoveData struct {
// 	Player    player        `json:"player"`
// 	Phase     gHeader.Phase `json:"phase"`
// 	Turn      int           `json:"turn"`
// 	card      card          `json:"card"`
// 	From      area          `json:"from"`
// 	ToHand    bool          `json:"toHand"`
// 	CreatedAt time.Time     `json:"createdAt"`
// 	Color     color.Color   `json:"color"`
// }
//
// func (g game) claimItemMoveData(p player, from area, cd card, toHand bool) claimItemMoveData {
// 	return claimItemMoveData{
// 		Player:    p.hideCards(),
// 		Phase:     g.Phase,
// 		Turn:      g.Turn,
// 		card:      cd,
// 		From:      from,
// 		ToHand:    toHand,
// 		CreatedAt: time.Now(),
// 	}
// }
//
// func (g game) claimItemMove(p player, from area, cd card, toHand bool) move.Move {
// 	return move.Move{
// 		Name: claimItemID,
// 		Data: g.claimItemMoveData(p, from, cd, toHand),
// 	}
// }

func (g game) finalClaim(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	g.Phase = phaseFinalClaim
	for row := rowA; row <= lastRowFor(g.NumPlayers); row++ {
		for col := col1; col <= col8; col++ {
			a, found := g.grid.area(row, col)
			if found {
				p, found := playerByID(a.thief.pid, g.players)
				if found {
					cd := a.card
					a.card = newCard(cdNone, cdFaceDown)
					a.thief.pid = pidNone
					p.DiscardPile = append([]card{cd}, p.DiscardPile...)
				}
			}
		}
	}
	for _, p := range g.players {
		p.collectCards()
	}
}
