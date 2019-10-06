package main

import (
	"github.com/SlothNinja/user"
)

func (g game) updateClickablesFor(u user.User) game {
	canClick := g.canClick(u)
	g.grid.each(func(a area) area {
		a.clickable = canClick(a)
		return a
	})
	return g
}

// canClick returns a function specialized by current game context to test whether a player can click on
// a particular area in the grid.  The main benefit is the function provides a closure around area computions,
// essentially caching the results.
func (g game) canClick(u user.User) func(area) bool {
	ff := func(a area) bool { return false }
	cp := g.currentPlayerFor(u)
	if cp.id == noPID || cp.performedAction {
		return ff
	}

	switch g.Phase {
	case phasePlaceThieves:
		return func(a area) bool { return a.thief.pid == noPID }
	case phaseSelectThief:
		return func(a area) bool { return a.thief.pid == cp.id }
	case phaseMoveThief:
		var toAreas []area
		from := g.grid.area(g.selectedAreaID.row, g.selectedAreaID.column)
		switch {
		case g.playedCard.kind == cdLamp || g.playedCard.kind == cdSLamp:
			toAreas = g.grid.lampAreas(from)
		case g.playedCard.kind == cdCamel || g.playedCard.kind == cdSCamel:
			toAreas = g.grid.camelAreas(from)
		case g.playedCard.kind == cdSword:
			toAreas = g.grid.swordAreasFor(cp, from)
		case g.playedCard.kind == cdCarpet:
			toAreas = g.grid.carpetAreas(from)
		case g.playedCard.kind == cdTurban && g.stepped == 0:
			toAreas = g.grid.turban0Areas(from)
		case g.playedCard.kind == cdTurban && g.stepped == 1:
			toAreas = g.turban1Areas()
		case g.playedCard.kind == cdCoins:
			toAreas = g.coinsAreas()
		}
		return func(a area) bool { return hasArea(toAreas, a) }
	}
	return ff
}

func (g grid) isLampMove(from, to area) bool {
	return hasArea(g.lampAreas(from), to)
}

func (g grid) lampAreas(from area) []area {
	var as []area

	if from == noArea {
		return as
	}

	to := g.lampW(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.lampE(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.lampN(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.lampS(from)
	if to != noArea {
		as = append(as, to)
	}

	return as
}

func (g grid) lampW(from area) area {
	var to1, to2 area

	for col := from.column - 1; col >= col1; col-- {
		to1 = g.area(from.row, col)
		if !canMove(to1) {
			break
		}
		to2 = to1
	}
	return to2
}

func (g grid) lampE(from area) area {
	var to1, to2 area

	for col := from.column + 1; col <= g.numCols(); col++ {
		to1 = g.area(from.row, col)
		if !canMove(to1) {
			break
		}
		to2 = to1
	}
	return to2
}

func (g grid) lampN(from area) area {
	var to1, to2 area

	for row := from.row - 1; row >= rowA; row-- {
		to1 = g.area(row, from.column)
		if !canMove(to1) {
			break
		}
		to2 = to1
	}
	return to2
}

func (g grid) lampS(from area) area {
	var to1, to2 area

	for row := from.row + 1; row <= g.numRows(); row++ {
		to1 = g.area(row, from.column)
		if !canMove(to1) {
			break
		}
		to2 = to1
	}
	return to2
}

func (g grid) isCamelMove(from, to area) bool {
	return hasArea(g.camelAreas(from), to)
}

func (g grid) camelAreas(from area) []area {
	var as []area

	if from == noArea {
		return as
	}

	to := g.camelWWW(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelEEE(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelNNN(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelSSS(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelWNW(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelWSW(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelENE(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelESE(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelSSE(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelNNE(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelSSW(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelNNW(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelN(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelE(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelS(from)
	if to != noArea {
		as = append(as, to)
	}

	to = g.camelW(from)
	if to != noArea {
		as = append(as, to)
	}

	return as
}

func (g grid) camelWWW(from area) area {
	// Final destination
	to1 := g.area(from.row, from.column-3)
	if !canMove(to1) {
		return noArea
	}

	to2 := g.area(from.row, from.column-1)
	to3 := g.area(from.row, from.column-2)

	if canMove(to2, to3) {
		return to1
	}

	return noArea
}

func (g grid) camelEEE(from area) area {
	// Final destination
	to1 := g.area(from.row, from.column+3)
	if !canMove(to1) {
		return noArea
	}

	to2 := g.area(from.row, from.column+1)
	to3 := g.area(from.row, from.column+2)

	if canMove(to2, to3) {
		return to1
	}

	return noArea
}

func (g grid) camelNNN(from area) area {
	// Final destination
	to1 := g.area(from.row-3, from.column)
	if !canMove(to1) {
		return noArea
	}

	to2 := g.area(from.row-1, from.column)
	to3 := g.area(from.row-2, from.column)

	if canMove(to2, to3) {
		return to1
	}

	return noArea
}

func (g grid) camelSSS(from area) area {
	// Final destination
	to1 := g.area(from.row+3, from.column)
	if !canMove(to1) {
		return noArea
	}

	to2 := g.area(from.row+1, from.column)
	to3 := g.area(from.row+2, from.column)

	if canMove(to2, to3) {
		return to1
	}

	return noArea
}

func (g grid) camelWNW(from area) area {
	// Final destination
	to1 := g.area(from.row-1, from.column-2)
	if !canMove(to1) {
		return noArea
	}

	// Try path of two left, one up
	to2 := g.area(from.row, from.column-1)
	to3 := g.area(from.row, from.column-2)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one left, one up, one left
	to4 := g.area(from.row-1, from.column-1)
	if canMove(to2, to4) {
		return to1
	}

	// Try path of one up, two left
	to5 := g.area(from.row-1, from.column)
	if canMove(to5, to4) {
		return to1
	}

	return noArea
}

func (g grid) camelWSW(from area) area {
	// Final destination
	to1 := g.area(from.row+1, from.column-2)
	if !canMove(to1) {
		return noArea
	}

	// Try path of two left, one down
	to2 := g.area(from.row, from.column-1)
	to3 := g.area(from.row, from.column-2)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one left, one down, one left
	to4 := g.area(from.row+1, from.column-1)
	if canMove(to2, to4) {
		return to1
	}

	// Try path of one down, two left
	to5 := g.area(from.row+1, from.column)
	if canMove(to5, to4) {
		return to1
	}

	return noArea
}

func (g grid) camelENE(from area) area {
	// Final destination
	to1 := g.area(from.row-1, from.column+2)
	if !canMove(to1) {
		return noArea
	}

	// Try path of two right, one up
	to2 := g.area(from.row, from.column+1)
	to3 := g.area(from.row, from.column+2)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one right, one up, one right
	to4 := g.area(from.row-1, from.column+1)
	if canMove(to2, to4) {
		return to1
	}

	// Try path of one up, two right
	to5 := g.area(from.row-1, from.column)
	if canMove(to5, to4) {
		return to1
	}

	return noArea
}

func (g grid) camelESE(from area) area {
	// Final destination
	to1 := g.area(from.row+1, from.column+2)
	if !canMove(to1) {
		return noArea
	}

	// Try path of two right, one down
	to2 := g.area(from.row, from.column+1)
	to3 := g.area(from.row, from.column+2)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one right, one down, one right
	to4 := g.area(from.row+1, from.column+1)
	if canMove(to2, to4) {
		return to1
	}

	// Try path of one up, two right
	to5 := g.area(from.row+1, from.column)
	if canMove(to5, to4) {
		return to1
	}

	return noArea
}

func (g grid) camelSSE(from area) area {
	// Final destination
	to1 := g.area(from.row+2, from.column+1)
	if !canMove(to1) {
		return noArea
	}

	// Try path of two down, one right
	to2 := g.area(from.row+1, from.column)
	to3 := g.area(from.row+2, from.column)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one down, one right, one down
	to4 := g.area(from.row+1, from.column+1)
	if canMove(to2, to4) {
		return to1
	}

	// Try path of one right, two down
	to5 := g.area(from.row, from.column+1)
	if canMove(to5, to4) {
		return to1
	}

	return noArea
}

func (g grid) camelNNE(from area) area {
	// Final destination
	to1 := g.area(from.row-2, from.column+1)
	if canMove(to1) {
		return noArea
	}

	// Try path of two up, one right
	to2 := g.area(from.row-1, from.column)
	to3 := g.area(from.row-2, from.column)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one up, one right, one up
	to4 := g.area(from.row-1, from.column+1)
	if canMove(to2, to4) {
		return to1
	}

	// Try path of one right, two up
	to5 := g.area(from.row, from.column+1)
	if canMove(to5, to4) {
		return to1
	}

	return noArea
}

func (g grid) camelSSW(from area) area {
	// Final destination
	to1 := g.area(from.row+2, from.column-1)
	if !canMove(to1) {
		return noArea
	}

	// Try path of two down, one left
	to2 := g.area(from.row+1, from.column)
	to3 := g.area(from.row+2, from.column)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one down, one left, one down
	to4 := g.area(from.row+1, from.column-1)
	if canMove(to2, to4) {
		return to1
	}

	// Try path of one left, two down
	to5 := g.area(from.row, from.column-1)
	if canMove(to5, to4) {
		return to1
	}

	return noArea
}

func (g grid) camelNNW(from area) area {
	// Final destination
	to1 := g.area(from.row-2, from.column-1)
	if !canMove(to1) {
		return noArea
	}

	// Try path of two up, one left
	to2 := g.area(from.row-1, from.column)
	to3 := g.area(from.row-2, from.column)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one up, one left, one up
	to4 := g.area(from.row-1, from.column-1)
	if canMove(to2, to4) {
		return to1
	}

	// Try path of one left, two up
	to5 := g.area(from.row, from.column-1)
	if canMove(to5, to4) {
		return to1
	}

	return noArea
}

func (g grid) camelN(from area) area {
	// Final destination
	to1 := g.area(from.row-1, from.column)
	if !canMove(to1) {
		return noArea
	}

	// Try path of one left, one up, one right
	to2 := g.area(from.row, from.column-1)
	to3 := g.area(from.row-1, from.column-1)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one right, one up, one left
	to4 := g.area(from.row, from.column+1)
	to5 := g.area(from.row-1, from.column+1)
	if canMove(to4, to5) {
		return to1
	}

	return noArea
}

func (g grid) camelE(from area) area {
	// Final destination
	to1 := g.area(from.row, from.column+1)
	if !canMove(to1) {
		return noArea
	}

	// Try path of one up, one right, one down
	to2 := g.area(from.row-1, from.column)
	to3 := g.area(from.row-1, from.column+1)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one down, one right, one up
	to4 := g.area(from.row-1, from.column)
	to5 := g.area(from.row-1, from.column+1)
	if canMove(to4, to5) {
		return to1
	}

	return noArea
}

func (g grid) camelS(from area) area {
	// Final destination
	to1 := g.area(from.row+1, from.column)
	if !canMove(to1) {
		return noArea
	}

	// Try path of one right, one down, one left
	to2 := g.area(from.row, from.column+1)
	to3 := g.area(from.row+1, from.column+1)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one left, one down, one right
	to4 := g.area(from.row, from.column-1)
	to5 := g.area(from.row+1, from.column-1)
	if canMove(to4, to5) {
		return to1
	}

	return noArea
}

func (g grid) camelW(from area) area {
	// Final destination
	to1 := g.area(from.row, from.column-1)
	if !canMove(to1) {
		return noArea
	}

	// Try path of one down, one left, one up
	to2 := g.area(from.row+1, from.column)
	to3 := g.area(from.row+1, from.column-1)
	if canMove(to2, to3) {
		return to1
	}

	// Try path of one up, one left, one down
	to4 := g.area(from.row-1, from.column)
	to5 := g.area(from.row-1, from.column-1)
	if canMove(to4, to5) {
		return to1
	}

	return noArea
}

func canMove(toAreas ...area) bool {
	if len(toAreas) == 0 {
		return false
	}

	for _, a := range toAreas {
		if a.hasThief() || !a.hasCard() {
			return false
		}
	}
	return true
}

func (g grid) isSwordMoveFor(cp player, from, to area) bool {
	return hasArea(g.swordAreasFor(cp, from), to)
}

func (g grid) swordAreasFor(cp player, from area) []area {
	var toAreas []area

	if from == noArea {
		return toAreas
	}

	to := g.swordWFor(cp, from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	to = g.swordEFor(cp, from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	to = g.swordNFor(cp, from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	to = g.swordSFor(cp, from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}
	return toAreas
}

func (g grid) swordWFor(cp player, from area) area {
	var to area
	for row, col := from.row, from.column-1; col >= col1; col-- {
		to = g.area(row, col)
		if !canMove(to) {
			break
		}
	}

	if !to.hasOtherThief(cp) {
		return noArea
	}

	bumpedTo := g.bumpedTo(from, to)
	if !canMove(bumpedTo) {
		return noArea
	}
	return to
}

func (g grid) swordEFor(cp player, from area) area {
	var to area
	for row, col := from.row, from.column+1; col <= g.numCols(); col++ {
		to = g.area(row, col)
		if !canMove(to) {
			break
		}
	}

	if !to.hasOtherThief(cp) {
		return noArea
	}

	bumpedTo := g.bumpedTo(from, to)
	if !canMove(bumpedTo) {
		return noArea
	}
	return to
}

func (g grid) swordNFor(cp player, from area) area {
	var to area
	for row, col := from.row-1, from.column; row >= rowA; row-- {
		to = g.area(row, col)
		if !canMove(to) {
			break
		}
	}

	if !to.hasOtherThief(cp) {
		return noArea
	}

	bumpedTo := g.bumpedTo(from, to)
	if !canMove(bumpedTo) {
		return noArea
	}
	return to
}

func (g grid) swordSFor(cp player, from area) area {
	var to area
	for row, col := from.row+1, from.column; row <= g.numRows(); row++ {
		to = g.area(row, col)
		if !canMove(to) {
			break
		}
	}

	if !to.hasOtherThief(cp) {
		return noArea
	}

	bumpedTo := g.bumpedTo(from, to)
	if !canMove(bumpedTo) {
		return noArea
	}
	return to
}

func (g grid) isCarpetMove(from, to area) bool {
	return hasArea(g.carpetAreas(from), to)
}

func (g grid) carpetAreas(from area) []area {
	var toAreas []area

	if from == noArea {
		return toAreas
	}

	to := g.carpetW(from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	to = g.carpetE(from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	to = g.carpetN(from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	to = g.carpetS(from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	return toAreas
}

func (g grid) carpetW(from area) area {
	var to1, to2 area

	for col := from.column - 1; col >= col1; col-- {
		to1 = g.area(from.row, col)
		if !canMove(to1) {
			break
		}
	}

	if to1 == noArea {
		return noArea
	}

	if to1.hasCard() {
		return noArea
	}

	for col := to1.column - 1; col >= col1; col-- {
		to2 = g.area(from.row, col)
		if to2.hasCard() {
			break
		}
	}

	if !canMove(to2) {
		return noArea
	}

	return to2
}

func (g grid) carpetE(from area) area {
	var to1, to2 area

	for col := from.column + 1; col <= g.numCols(); col++ {
		to1 = g.area(from.row, col)
		if !canMove(to1) {
			break
		}
	}

	if to1 == noArea {
		return noArea
	}

	if to1.hasCard() {
		return noArea
	}

	for col := to1.column + 1; col <= g.numCols(); col++ {
		to2 = g.area(from.row, col)
		if to2.hasCard() {
			break
		}
	}

	if !canMove(to2) {
		return noArea
	}

	return to2
}

func (g grid) carpetS(from area) area {
	var to1, to2 area

	for row := from.row + 1; row <= g.numRows(); row++ {
		to1 = g.area(row, from.column)
		if !canMove(to1) {
			break
		}
	}

	if to1 == noArea {
		return noArea
	}

	if to1.hasCard() {
		return noArea
	}

	for row := to1.row + 1; row <= g.numRows(); row++ {
		to2 = g.area(row, from.column)
		if to2.hasCard() {
			break
		}
	}

	if !canMove(to2) {
		return noArea
	}

	return to2
}

func (g grid) carpetN(from area) area {
	var to1, to2 area

	for row := from.row - 1; row >= rowA; row-- {
		to1 = g.area(row, from.column)
		if !canMove(to1) {
			break
		}
	}

	if to1 == noArea {
		return noArea
	}

	if to1.hasCard() {
		return noArea
	}

	for row := to1.row - 1; row >= rowA; row-- {
		to2 = g.area(row, from.column)
		if to2.hasCard() {
			break
		}
	}

	if !canMove(to2) {
		return noArea
	}

	return to2
}

func (g grid) isTurban0Move(from, to area) bool {
	return hasArea(g.turban0Areas(from), to)
}

func (g grid) turban0Areas(from area) []area {
	var toAreas []area

	if from == noArea {
		return toAreas
	}

	if !from.hasThief() {
		return toAreas
	}

	to := g.turban0W(from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	to = g.turban0E(from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	to = g.turban0N(from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	to = g.turban0S(from)
	if to != noArea {
		toAreas = append(toAreas, to)
	}

	return toAreas
}

func (g grid) turban0W(from area) area {
	to := g.area(from.row, from.column-1)
	if !canMove(to) {
		return noArea
	}

	if !g.canMoveOrthogonal(to) {
		return noArea
	}
	return to
}

func (g grid) turban0E(from area) area {
	to := g.area(from.row, from.column+1)
	if !canMove(to) {
		return noArea
	}

	if !g.canMoveOrthogonal(to) {
		return noArea
	}
	return to
}

func (g grid) turban0N(from area) area {
	to := g.area(from.row-1, from.column)
	if !canMove(to) {
		return noArea
	}

	if !g.canMoveOrthogonal(to) {
		return noArea
	}
	return to
}

func (g grid) turban0S(from area) area {
	to := g.area(from.row+1, from.column)
	if !canMove(to) {
		return noArea
	}

	if !g.canMoveOrthogonal(to) {
		return noArea
	}
	return to
}

func (g grid) canMoveOrthogonal(from area) bool {
	toE := g.area(from.row, from.column+1)
	if canMove(toE) {
		return true
	}

	toW := g.area(from.row, from.column-1)
	if canMove(toW) {
		return true
	}

	toN := g.area(from.row-1, from.column)
	if canMove(toN) {
		return true
	}

	toS := g.area(from.row+1, from.column)
	return canMove(toS)
}

func (g game) isTurban1Area(a area) bool {
	return hasArea(g.turban1Areas(), a)
}

func (g game) turban1Areas() []area {
	var as []area

	a := g.grid.area(g.selectedAreaID.row, g.selectedAreaID.column)
	if a == noArea {
		return as
	}

	// Move Left
	a2 := g.grid.area(a.row, a.column-1)
	if canMove(a2) {
		as = append(as, a2)
	}

	// Move Right
	a2 = g.grid.area(a.row, a.column+1)
	if canMove(a2) {
		as = append(as, a2)
	}

	// Move Up
	a2 = g.grid.area(a.row-1, a.column)
	if canMove(a2) {
		as = append(as, a2)
	}

	// Move Down
	a2 = g.grid.area(a.row+1, a.column)
	if canMove(a2) {
		as = append(as, a2)
	}

	return as
}

func (g game) isCoinsArea(a area) bool {
	return hasArea(g.coinsAreas(), a)
}

func (g game) coinsAreas() []area {
	return g.turban1Areas()
}
