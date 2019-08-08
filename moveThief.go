package main

import (
	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// func init() {
// 	move.Register(moveThiefID, new(moveThiefMoveData))
// 	move.Register(bumpedThiefID, new(bumpedThiefMoveData))
// }

const moveThiefID = "move-thief"
const bumpedThiefID = "bumped-thief"

func (g game) startMoveThief() game {
	g.Phase = phaseMoveThief
	return g
}

// type moveThiefMoveData struct {
// 	Player    player        `json:"player"`
// 	Phase     gHeader.Phase `json:"phase"`
// 	Turn      int           `json:"turn"`
// 	From      area          `json:"from"`
// 	To        area          `json:"to"`
// 	CreatedAt time.Time     `json:"createdAt"`
// 	Color     color.Color   `json:"color"`
// }

// func (g game) moveThiefMoveData(p player, from, to area) moveThiefMoveData {
// 	return moveThiefMoveData{
// 		Player:    p.hideCards(),
// 		Phase:     g.Phase,
// 		Turn:      g.Turn,
// 		From:      from,
// 		To:        to,
// 		CreatedAt: time.Now(),
// 	}
// }

// func (g game) moveThiefMove(p player, from, to area) move.Move {
// 	return move.Move{
// 		Name: moveThiefID,
// 		Data: g.moveThiefMoveData(p, from, to),
// 	}
// }

// func (m moveThiefMoveData) colorize(g game, u user.User2) {
// 	m.Color = g.colorByPIDFor(u)(m.Player.ID)
// }

// type bumpedThiefMoveData struct {
// 	Player    player        `json:"player"`
// 	Phase     gHeader.Phase `json:"phase"`
// 	Turn      int           `json:"turn"`
// 	From      area          `json:"from"`
// 	To        area          `json:"to"`
// 	CreatedAt time.Time     `json:"createdAt"`
// 	Color     color.Color   `json:"color"`
// }
//
// func (g game) bumpedThiefMoveData(p player, from, to area) bumpedThiefMoveData {
// 	return bumpedThiefMoveData{
// 		Player:    p.hideCards(),
// 		Phase:     g.Phase,
// 		Turn:      g.Turn,
// 		From:      from,
// 		To:        to,
// 		CreatedAt: time.Now(),
// 	}
// }

// func (g game) bumpedThiefMove(p player, from, to area) move.Move {
// 	return move.Move{
// 		Name: bumpedThiefID,
// 		Data: g.bumpedThiefMoveData(p, from, to),
// 	}
// }

// func (b bumpedThiefMoveData) colorize(g game, u user.User2) {
// 	b.Color = g.colorByPIDFor(u)(b.Player.ID)
// }

func (g game) MoveThief(c *gin.Context) (game, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, from, to, cd, err := g.validateMoveThief(c)
	if err != nil {
		return g, err
	}

	g.Log = append(g.Log, logEntry{
		"template": moveThiefID,
		"pid":      cp.ID,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"card":     cd,
		"from":     from,
		"to":       to,
	})
	// Log move thief
	// g.GLog.AddEntryData(glog.EntryData{
	// 	"template": moveThiefID,
	// 	"turn":     g.Turn,
	// 	"phase":    g.Phase,
	// 	"pid":      cp.ID,
	// 	"card":     cd,
	// 	"from":     from,
	// 	"to":       to,
	// })

	cp.Score += to.card.value()

	switch cd.kind {
	case cdSword:
		g.swordMove(cp, from, to)
	case cdTurban:
		g.turbanMove(cp, from, to)
	case cdCoins:
		g.coinMove(cp, from, to)
	default:
		g = g.defaultMove(cp, from, to, g.toHand())
	}

	cu, _ := user.Current(c)
	g = g.updateClickablesFor(cu)
	g.Stack = g.Stack.Update()
	return g, nil
}

func (g game) swordMove(cp player, from, to area) game {
	bumpedTo, found := g.bumpedTo(from, to)
	if found {
		bumpedPID := to.thief.pid
		bumpedTo.thief.pid = bumpedPID
		p2, found := playerByID(bumpedPID, g.players)
		if found {
			p2.Score += bumpedTo.card.value() - to.card.value()
		}
	}

	// Capture move data before updating state.
	// m := g.moveThiefMove(*cp, from, to)
	// g.Animations, g.Log = append(g.Animations, m), append(g.Log, m)

	// Move thief
	from.thief.pid, to.thief.pid = pidNone, cp.ID

	// Bump thief move
	// bm := g.bumpedThiefMove(*p2, to, *bumpedTo)
	//	g.Animations, g.Log = append(g.Animations, bm), append(g.Log, bm)

	// Claim Item
	g = g.claimItem(from, cp)
	cp.PerformedAction = true
	return g.updatePlayer(cp)
}

func (g game) turbanMove(cp player, from, to area) (game, player) {
	if g.stepped == 0 {
		g.stepped = 1
		g.selectedAreaID = to.areaID
		g = g.startMoveThief()

		g.defaultMove(cp, from, to, false)

		// Revised defaultMove
		cp.PerformedAction = false
		g.Phase = phaseMoveThief
	} else {
		g.stepped = 2
		g.defaultMove(cp, from, to, true)
	}
	return g, cp
}

func (g game) coinMove(cp player, from, to area) game {
	g = g.defaultMove(cp, from, to, true)
	cp, _, _ = cp.draw()
	return g.updatePlayer(cp)
}

func (g game) removeThiefFrom(a area) (game, area) {
	a.thief.pid = pidNone
	return g.updateArea(a), a
}

func (g game) moveThief(p player, from, to area) (game, area, area) {
	g, from = g.removeThiefFrom(from)
	to.thief.pid, to.thief.from = p.ID, from.areaID
	return g.updateArea(to), from, to
}

func (g game) defaultMove(cp player, from, to area, toHand bool) game {
	// Move thief
	g, from, to = g.moveThief(cp, from, to)

	// Claim Item
	g = g.claimItem(from, cp)
	if !toHand {
		cp, _, _ = cp.draw()
	}

	cp.PerformedAction = true
	return g.updatePlayer(cp)
}

func (g game) validateMoveThief(c *gin.Context) (player, area, area, card, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, err := g.validatePlayerAction(c)
	if err != nil {
		return player{}, area{}, area{}, card{}, err
	}

	to, err := g.getArea(c)
	if err != nil {
		return player{}, area{}, area{}, card{}, err
	}

	from, found := g.SelectedThiefArea()
	if !found {
		return player{}, area{}, area{}, card{}, errors.Wrap(errValidation, "selected thief area not found")
	}

	cd := g.playedCard
	switch {
	case from.thief.pid != cp.ID:
		return player{}, area{}, area{}, card{},
			errors.WithMessage(errValidation, "selected thief of another player")
	case cd.kind == cdNone:
		return player{}, area{}, area{}, card{},
			errors.WithMessage(errValidation, "you must play card before moving thief")
	case (cd.kind == cdLamp || cd.kind == cdSLamp) && !g.isLampArea(to),
		(cd.kind == cdCamel || cd.kind == cdSCamel) && !g.isCamelArea(to),
		cd.kind == cdCoins && !g.isCoinsArea(to),
		cd.kind == cdSword && !g.isSwordAreaFor(cp, to),
		cd.kind == cdCarpet && !g.isCarpetArea(to),
		cd.kind == cdTurban && g.stepped == 0 && !g.isTurban0Area(to),
		cd.kind == cdTurban && g.stepped == 1 && !g.isTurban1Area(to),
		cd.kind == cdGuard:
		return player{}, area{}, area{}, card{},
			errors.WithMessage(errValidation, "played card does not permit moving selected thief to selected area")
	}
	return cp, from, to, cd, nil
}

func (g game) bumpedTo(from, to area) (area, bool) {
	switch {
	case from.row > to.row:
		return g.grid.area(to.row-1, from.column)
	case from.row < to.row:
		return g.grid.area(to.row+1, from.column)
	case from.column > to.column:
		return g.grid.area(from.row, to.column-1)
	case from.column < to.column:
		return g.grid.area(from.row, to.column+1)
	}
	return area{}, false
}
