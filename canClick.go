package main

import (
	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"
)

func (g game) updateClickablesFor(u user.User2) game {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

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
func (g game) canClick(u user.User2) func(area) bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	ff := func(a area) bool { return false }
	cp, found := g.currentPlayerFor(u)
	if !found || cp.performedAction {
		return ff
	}

	switch g.Phase {
	case phasePlaceThieves:
		return func(a area) bool { return a.thief.pid == pidNone }
	case phaseSelectThief:
		return func(a area) bool { return a.thief.pid == cp.id }
	case phaseMoveThief:
		var as []area
		switch {
		case g.playedCard.kind == cdLamp || g.playedCard.kind == cdSLamp:
			as = g.lampAreas()
		case g.playedCard.kind == cdCamel || g.playedCard.kind == cdSCamel:
			as = g.camelAreas()
		case g.playedCard.kind == cdSword:
			as = g.swordAreasFor(cp)
		case g.playedCard.kind == cdCarpet:
			as = g.carpetAreas()
		case g.playedCard.kind == cdTurban && g.stepped == 0:
			as = g.turban0Areas()
		case g.playedCard.kind == cdTurban && g.stepped == 1:
			as = g.turban1Areas()
		case g.playedCard.kind == cdCoins:
			as = g.coinsAreas()
		}
		return func(a area) bool { return hasArea(as, a) }
	}
	return ff
}

func (g game) isLampArea(a area) bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return hasArea(g.lampAreas(), a)
}

func (g game) lampAreas() []area {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	var (
		as []area
		a2 area
	)

	a1, found := g.SelectedThiefArea()
	if !found {
		return as
	}

	// Move Left
	add := false
	for col := a1.column - 1; col >= col1; col-- {
		temp, found := g.grid.area(a1.row, col)
		if !found || !canMoveTo(temp) {
			break
		}
		a2, add = temp, true
	}
	if add {
		as = append(as, a2)
	}

	// Move right
	add = false
	for col := a1.column + 1; col <= g.grid.numCols(); col++ {
		temp, found := g.grid.area(a1.row, col)
		if !found || !canMoveTo(temp) {
			break
		}
		a2, add = temp, true
	}
	if add {
		as = append(as, a2)
	}

	// Move Up
	add = false
	for row := a1.row - 1; row >= rowA; row-- {
		temp, found := g.grid.area(row, a1.column)
		if !found || !canMoveTo(temp) {
			break
		}
		a2, add = temp, true
	}

	if add {
		as = append(as, a2)
	}

	// Move Down
	add = false
	for row := a1.row + 1; row <= g.grid.numRows(); row++ {
		temp, found := g.grid.area(row, a1.column)
		if !found || !canMoveTo(temp) {
			break
		}
		a2, add = temp, true
	}
	if a2.row != rowNone {
		as = append(as, a2)
	}

	return as
}

func (g game) isCamelArea(a area) bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return hasArea(g.camelAreas(), a)
}

func (g game) camelAreas() []area {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	var as []area

	a, found := g.SelectedThiefArea()
	if !found {
		return as
	}

	// Move Three Left?
	if a.column-3 >= col1 {
		area1, found1 := g.grid.area(a.row, a.column-1)
		area2, found2 := g.grid.area(a.row, a.column-2)
		area3, found3 := g.grid.area(a.row, a.column-3)
		if found1 && found2 && found3 && canMoveTo(area1, area2, area3) {
			as = append(as, area3)
		}
	}

	// Move Three Right?
	if a.column+3 <= g.grid.numCols() {
		area1, found1 := g.grid.area(a.row, a.column+1)
		area2, found2 := g.grid.area(a.row, a.column+2)
		area3, found3 := g.grid.area(a.row, a.column+3)
		if found1 && found2 && found3 && canMoveTo(area1, area2, area3) {
			as = append(as, area3)
		}
	}

	// Move Three Up?
	if a.row-3 >= rowA {
		area1, found1 := g.grid.area(a.row-1, a.column)
		area2, found2 := g.grid.area(a.row-2, a.column)
		area3, found3 := g.grid.area(a.row-3, a.column)
		if found1 && found2 && found3 && canMoveTo(area1, area2, area3) {
			as = append(as, area3)
		}
	}

	// Move Three Down?
	if a.row+3 <= g.grid.numRows() {
		area1, found1 := g.grid.area(a.row+1, a.column)
		area2, found2 := g.grid.area(a.row+2, a.column)
		area3, found3 := g.grid.area(a.row+3, a.column)
		if found1 && found2 && found3 && canMoveTo(area1, area2, area3) {
			as = append(as, area3)
		}
	}

	// Move Two Left One Up or One Up Two Left or One Left One Up One Left?
	if a.column-2 >= col1 && a.row-1 >= rowA {
		area1, found1 := g.grid.area(a.row, a.column-1)
		area2, found2 := g.grid.area(a.row, a.column-2)
		area3, found3 := g.grid.area(a.row-1, a.column-2)
		area4, found4 := g.grid.area(a.row-1, a.column)
		area5, found5 := g.grid.area(a.row-1, a.column-1)
		if (found1 && found2 && found3 && canMoveTo(area1, area2, area3)) ||
			(found3 && found4 && found5 && canMoveTo(area3, area4, area5)) ||
			(found1 && found5 && found3 && canMoveTo(area1, area5, area3)) {
			as = append(as, area3)
		}
	}

	// Move Two Left One Down or One Down Two Left or One Left One Down One Left?
	if a.column-2 >= col1 && a.row+1 <= g.grid.numRows() {
		area1, found1 := g.grid.area(a.row, a.column-1)
		area2, found2 := g.grid.area(a.row, a.column-2)
		area3, found3 := g.grid.area(a.row+1, a.column-2)
		area4, found4 := g.grid.area(a.row+1, a.column)
		area5, found5 := g.grid.area(a.row+1, a.column-1)
		if (found1 && found2 && found3 && canMoveTo(area1, area2, area3)) ||
			(found3 && found4 && found5 && canMoveTo(area3, area4, area5)) ||
			(found1 && found5 && found3 && canMoveTo(area1, area5, area3)) {
			as = append(as, area3)
		}
	}

	// Move Two Right One Up or One Up Two Right or One Right One Up One Right?
	if a.column+2 <= g.grid.numCols() && a.row-1 >= rowA {
		area1, found1 := g.grid.area(a.row, a.column+1)
		area2, found2 := g.grid.area(a.row, a.column+2)
		area3, found3 := g.grid.area(a.row-1, a.column+2)
		area4, found4 := g.grid.area(a.row-1, a.column)
		area5, found5 := g.grid.area(a.row-1, a.column+1)
		if (found1 && found2 && found3 && canMoveTo(area1, area2, area3)) ||
			(found3 && found4 && found5 && canMoveTo(area3, area4, area5)) ||
			(found1 && found5 && found3 && canMoveTo(area1, area5, area3)) {
			as = append(as, area3)
		}
	}

	// Move Two Right One Down or One Down Two Right or One Right One Down One Right?
	if a.column+2 <= g.grid.numCols() && a.row+1 <= g.grid.numRows() {
		area1, found1 := g.grid.area(a.row, a.column+1)
		area2, found2 := g.grid.area(a.row, a.column+2)
		area3, found3 := g.grid.area(a.row+1, a.column+2)
		area4, found4 := g.grid.area(a.row+1, a.column)
		area5, found5 := g.grid.area(a.row+1, a.column+1)
		if (found1 && found2 && found3 && canMoveTo(area1, area2, area3)) ||
			(found3 && found4 && found5 && canMoveTo(area3, area4, area5)) ||
			(found1 && found5 && found3 && canMoveTo(area1, area5, area3)) {
			as = append(as, area3)
		}
	}

	// Move One Right Two Down or Two Down One Right or One Down One Right One Down?
	if a.column+1 <= g.grid.numCols() && a.row+2 <= g.grid.numRows() {
		area1, found1 := g.grid.area(a.row+1, a.column)
		area2, found2 := g.grid.area(a.row+2, a.column)
		area3, found3 := g.grid.area(a.row+2, a.column+1)
		area4, found4 := g.grid.area(a.row, a.column+1)
		area5, found5 := g.grid.area(a.row+1, a.column+1)
		if (found1 && found2 && found3 && canMoveTo(area1, area2, area3)) ||
			(found3 && found4 && found5 && canMoveTo(area3, area4, area5)) ||
			(found1 && found5 && found3 && canMoveTo(area1, area5, area3)) {
			as = append(as, area3)
		}
	}

	// Move One Right Two Up or Two Up One Right or One Up One Right One Up?
	if a.column+1 <= g.grid.numCols() && a.row-2 >= rowA {
		area1, found1 := g.grid.area(a.row-1, a.column)
		area2, found2 := g.grid.area(a.row-2, a.column)
		area3, found3 := g.grid.area(a.row-2, a.column+1)
		area4, found4 := g.grid.area(a.row, a.column+1)
		area5, found5 := g.grid.area(a.row-1, a.column+1)
		if (found1 && found2 && found3 && canMoveTo(area1, area2, area3)) ||
			(found3 && found4 && found5 && canMoveTo(area3, area4, area5)) ||
			(found1 && found5 && found3 && canMoveTo(area1, area5, area3)) {
			as = append(as, area3)
		}
	}

	// Move One Left Two Down or Two Down One Left or One Down One Left One Down?
	if a.column-1 >= col1 && a.row+2 <= g.grid.numRows() {
		area1, found1 := g.grid.area(a.row+1, a.column)
		area2, found2 := g.grid.area(a.row+2, a.column)
		area3, found3 := g.grid.area(a.row+2, a.column-1)
		area4, found4 := g.grid.area(a.row, a.column-1)
		area5, found5 := g.grid.area(a.row+1, a.column-1)
		if (found1 && found2 && found3 && canMoveTo(area1, area2, area3)) ||
			(found3 && found4 && found5 && canMoveTo(area3, area4, area5)) ||
			(found1 && found5 && found3 && canMoveTo(area1, area5, area3)) {
			as = append(as, area3)
		}
	}

	// Move One Left Two Up or Two Up One Left or One Up One Left One Up?
	if a.column-1 >= col1 && a.row-2 >= rowA {
		area1, found1 := g.grid.area(a.row-1, a.column)
		area2, found2 := g.grid.area(a.row-2, a.column)
		area3, found3 := g.grid.area(a.row-2, a.column-1)
		area4, found4 := g.grid.area(a.row, a.column-1)
		area5, found5 := g.grid.area(a.row-1, a.column-1)
		if (found1 && found2 && found3 && canMoveTo(area1, area2, area3)) ||
			(found3 && found4 && found5 && canMoveTo(area3, area4, area5)) ||
			(found1 && found5 && found3 && canMoveTo(area1, area5, area3)) {
			as = append(as, area3)
		}
	}

	// Move One Left One Up One Right or One Up One Left One Down?
	if a.column-1 >= col1 && a.row-1 >= rowA {
		area1, found1 := g.grid.area(a.row, a.column-1)
		area2, found2 := g.grid.area(a.row-1, a.column-1)
		area3, found3 := g.grid.area(a.row-1, a.column)
		if found1 && found2 && found3 && canMoveTo(area1, area2, area3) {
			as = append(as, area1, area3)
		}
	}

	// Move One Up One Right One Down or One Right One Up One Left?
	if a.column+1 <= g.grid.numCols() && a.row-1 >= rowA {
		area1, found1 := g.grid.area(a.row, a.column+1)
		area2, found2 := g.grid.area(a.row-1, a.column+1)
		area3, found3 := g.grid.area(a.row-1, a.column)
		if found1 && found2 && found3 && canMoveTo(area1, area2, area3) {
			as = append(as, area1, area3)
		}
	}

	// Move One Left One Down One Right or One Down One Left One Up?
	if a.column-1 >= col1 && a.row+1 <= g.grid.numRows() {
		area1, found1 := g.grid.area(a.row, a.column-1)
		area2, found2 := g.grid.area(a.row+1, a.column-1)
		area3, found3 := g.grid.area(a.row+1, a.column)
		if found1 && found2 && found3 && canMoveTo(area1, area2, area3) {
			as = append(as, area1, area3)
		}
	}

	// Move One Down One Right One Up or One Right One Down One Left?
	if a.column+1 <= g.grid.numCols() && a.row+1 <= g.grid.numRows() {
		area1, found1 := g.grid.area(a.row, a.column+1)
		area2, found2 := g.grid.area(a.row+1, a.column+1)
		area3, found3 := g.grid.area(a.row+1, a.column)
		if found1 && found2 && found3 && canMoveTo(area1, area2, area3) {
			as = append(as, area1, area3)
		}
	}

	return as
}

func canMoveTo(as ...area) bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	for _, a := range as {
		if a.hasThief() || !a.hasCard() {
			return false
		}
	}
	return true
}

func (g game) isSwordAreaFor(cp player, a area) bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return hasArea(g.swordAreasFor(cp), a)
}

func (g game) swordAreasFor(cp player) []area {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	var as []area

	a, found := g.SelectedThiefArea()
	if !found {
		return as
	}

	///////////////////////////////////////
	// Move Left

	// Left as far as permitted
	for row, col := a.row, a.column-1; col >= col1; col-- {
		temp, found := g.grid.area(row, col)
		if !found {
			break
		}
		if canMoveTo(temp) {
			continue
		}
		if temp.hasOtherThief(cp) {
			bumpTo, found2 := g.grid.area(row, col-1)
			if found2 && canMoveTo(bumpTo) {
				as = append(as, temp)
			}
			break
		}
	}

	/////////////////////////////////////////////
	// Move Right

	// Right as far as permitted
	for row, col := a.row, a.column+1; col <= g.grid.numCols(); col++ {
		temp, found := g.grid.area(row, col)
		if !found {
			break
		}
		if canMoveTo(temp) {
			continue
		}
		if temp.hasOtherThief(cp) {
			bumpTo, found2 := g.grid.area(row, col+1)
			if found2 && canMoveTo(bumpTo) {
				as = append(as, temp)
			}
			break
		}
	}

	//////////////////////////////////////////////////
	// Move Up

	// Up as far as permitted
	for row, col := a.row-1, a.column; row >= rowA; row-- {
		temp, found := g.grid.area(row, col)
		if !found {
			break
		}
		if canMoveTo(temp) {
			continue
		}
		if temp.hasOtherThief(cp) {
			bumpTo, found2 := g.grid.area(row-1, col)
			if found2 && canMoveTo(bumpTo) {
				as = append(as, temp)
			}
			break
		}
	}

	////////////////////////////////////
	// Move Down

	// Down as far as permitted
	for row, col := a.row+1, a.column; row <= g.grid.numRows(); row++ {
		temp, found := g.grid.area(row, col)
		if !found {
			break
		}
		if canMoveTo(temp) {
			continue
		}
		if temp.hasOtherThief(cp) {
			bumpTo, found2 := g.grid.area(row+1, col)
			if found2 && canMoveTo(bumpTo) {
				as = append(as, temp)
			}
			break
		}
	}

	return as
}

func (g game) isCarpetArea(a area) bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return hasArea(g.carpetAreas(), a)
}

func (g game) carpetAreas() []area {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	var (
		as []area
		a2 area
	)

	a1, found := g.SelectedThiefArea()
	if !found {
		return as
	}

	///////////////////////////////////
	// Move Left
	found1, add := false, false

MoveLeft:
	for col := a1.column - 1; col >= col1; col-- {
		switch temp, found := g.grid.area(a1.row, col); {
		case found && !temp.hasCard():
			found1 = true
		case found1 && canMoveTo(temp):
			a2, add = temp, true
			break MoveLeft
		default:
			break MoveLeft
		}
	}
	if add {
		as = append(as, a2)
	}

	////////////////////////////////////////////
	// Move Right
	found1, add = false, false

MoveRight:
	for col := a1.column + 1; col <= g.grid.numCols(); col++ {
		switch temp, found := g.grid.area(a1.row, col); {
		case found && !temp.hasCard():
			found1 = true
		case found1 && canMoveTo(temp):
			a2, add = temp, true
			break MoveRight
		default:
			break MoveRight
		}
	}
	if add {
		as = append(as, a2)
	}

	/////////////////////////////////////////
	// Move Up
	found1, add = false, false

MoveUp:
	for row := a1.row - 1; row >= rowA; row-- {
		switch temp, found := g.grid.area(row, a1.column); {
		case found && !temp.hasCard():
			found1 = true
		case found1 && canMoveTo(temp):
			a2, add = temp, true
			break MoveUp
		default:
			break MoveUp
		}
	}
	if add {
		as = append(as, a2)
	}

	////////////////////////////////////////////////
	// Move Down
	found1, add = false, false

MoveDown:
	for row := a1.row + 1; row <= g.grid.numRows(); row++ {
		switch temp, found := g.grid.area(row, a1.column); {
		case found && temp.hasCard():
			found1 = true
		case found1 && canMoveTo(temp):
			a2, add = temp, true
			break MoveDown
		default:
			break MoveDown
		}
	}
	if add {
		as = append(as, a2)
	}

	return as
}

func (g game) isTurban0Area(a area) bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return hasArea(g.turban0Areas(), a)
}

func (g game) turban0Areas() []area {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	var (
		as []area
		a2 area
	)

	a, found := g.SelectedThiefArea()
	if !found {
		return as
	}

	// Move Left
	a2, found2 := g.grid.area(a.row, a.column-1)
	if found2 && canMoveTo(a2) {
		// Left
		a3, found3 := g.grid.area(a2.row, a2.column-1)

		// Up
		a4, found4 := g.grid.area(a2.row-1, a2.column)

		// Down
		a5, found5 := g.grid.area(a2.row+1, a2.column)

		// Add
		if (found3 && canMoveTo(a3)) ||
			(found4 && canMoveTo(a4)) ||
			(found5 && canMoveTo(a5)) {
			as = append(as, a2)
		}
	}

	// Move Right
	a2, found2 = g.grid.area(a.row, a.column+1)
	if found2 && canMoveTo(a2) {
		// Right
		a3, found3 := g.grid.area(a2.row, a2.column+1)

		// Up
		a4, found4 := g.grid.area(a2.row-1, a2.column)

		// Down
		a5, found5 := g.grid.area(a2.row+1, a2.column)

		// Add
		if (found3 && canMoveTo(a3)) ||
			(found4 && canMoveTo(a4)) ||
			(found5 && canMoveTo(a5)) {
			as = append(as, a2)
		}
	}

	// Move Up
	a2, found2 = g.grid.area(a.row-1, a.column)
	if found2 && canMoveTo(a2) {
		// Left
		a3, found3 := g.grid.area(a2.row, a2.column-1)

		// Right
		a4, found4 := g.grid.area(a2.row, a2.column+1)

		// Up
		a5, found5 := g.grid.area(a2.row-1, a2.column)

		// Add
		if (found3 && canMoveTo(a3)) ||
			(found4 && canMoveTo(a4)) ||
			(found5 && canMoveTo(a5)) {
			as = append(as, a2)
		}
	}

	// Move Down
	a2, found2 = g.grid.area(a.row+1, a.column)
	if found2 && canMoveTo(a2) {
		// Left
		a3, found3 := g.grid.area(a2.row, a2.column-1)

		// Right
		a4, found4 := g.grid.area(a2.row, a2.column+1)

		// Down
		a5, found5 := g.grid.area(a2.row+1, a2.column)

		// Add
		if (found3 && canMoveTo(a3)) ||
			(found4 && canMoveTo(a4)) ||
			(found5 && canMoveTo(a5)) {
			as = append(as, a2)
		}
	}

	return as
}

func (g game) isTurban1Area(a area) bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return hasArea(g.turban1Areas(), a)
}

func (g game) turban1Areas() []area {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	var as []area

	a, found := g.SelectedThiefArea()
	if !found {
		return as
	}

	// Move Left
	a2, found := g.grid.area(a.row, a.column-1)
	if found && canMoveTo(a2) {
		as = append(as, a2)
	}

	// Move Right
	a2, found = g.grid.area(a.row, a.column+1)
	if found && canMoveTo(a2) {
		as = append(as, a2)
	}

	// Move Up
	a2, found = g.grid.area(a.row-1, a.column)
	if found && canMoveTo(a2) {
		as = append(as, a2)
	}

	// Move Down
	a2, found = g.grid.area(a.row+1, a.column)
	if found && canMoveTo(a2) {
		as = append(as, a2)
	}

	return as
}

func (g game) isCoinsArea(a area) bool {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return hasArea(g.coinsAreas(), a)
}

func (g game) coinsAreas() []area {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return g.turban1Areas()
}
