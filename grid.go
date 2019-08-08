package main

import (
	"encoding/json"
	"math/rand"
	"strconv"
)

type grid [][]area

const (
	rowNone int = iota
	rowA
	rowB
	rowC
	rowD
	rowE
	rowF
	rowG
)

var rowIDStrings = map[int]string{rowNone: "None", rowA: "A", rowB: "B", rowC: "C",
	rowD: "D", rowE: "E", rowF: "F", rowG: "G"}

// rowString outputs a row label.
func rowString(row int) string {
	return rowIDStrings[row]
}

// rowIDString outputs an row id.
func rowIDString(row int) string {
	return strconv.Itoa(row)
}

const (
	colNone int = iota
	col1
	col2
	col3
	col4
	col5
	col6
	col7
	col8
)

var columnIDStrings = map[int]string{colNone: "None", col1: "1", col2: "2", col3: "3", col4: "4",
	col5: "5", col6: "6", col7: "7", col8: "8"}

// ColString outputs a column label.
func ColString(col int) string {
	return columnIDStrings[col]
}

// ColIDString outputs an column id.
func ColIDString(col int) string {
	return strconv.Itoa(col)
}

func lastRowFor(numPlayers int) (row int) {
	row = rowG
	if numPlayers == 2 {
		row = rowF
	}
	return
}

func newGrid(numPlayers int) grid {
	deck := newDeck()
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	lastRow := lastRowFor(numPlayers)
	grid := make([][]area, lastRow)
	for row := 0; row < lastRow; row++ {
		grid[row] = make([]area, col8)
		for col := 0; col < col8; col++ {
			var cd card
			deck, cd = draw(deck)
			cd.turn(cdFaceUp)
			grid[row][col] = newArea(row+1, col+1, cd)
		}
	}
	return grid
}

func (g grid) area(row, col int) (area, bool) {
	found := row >= 1 && col >= 1 && row <= g.NumRows() && col <= g.NumCols()
	if found {
		return g[row-1][col-1], true
	}
	return area{}, false
}

func (g grid) NumRows() int {
	return len(g)
}

func (g grid) NumCols() int {
	if g.NumRows() > 1 {
		return len(g[0])
	}
	return 0
}

func (g grid) Each(f func(a area) area) {
	for row := range g {
		for col := range g[row] {
			g[row][col] = f(g[row][col])
		}
	}
}

func (g game) updateArea(a area) game {
	g.grid[a.row-1][a.column-1] = a
	return g
}

func (g grid) MarshalJSON() ([]byte, error) {
	type jGrid grid
	return json.Marshal(jGrid(g))
}

func (g *grid) UnmarshalJSON(v []byte) error {
	var unmarshaled [][]area
	err := json.Unmarshal(v, &unmarshaled)
	if err == nil {
		*g = unmarshaled
	}
	return err
}

type area struct {
	areaID
	thief     thief
	card      card
	clickable bool
}

type jArea struct {
	jAreaID
	Thief     thief `json:"thief"`
	Card      card  `json:"card"`
	Clickable bool  `json:"clickable"`
}

func (a area) MarshalJSON() ([]byte, error) {
	j := jArea{
		Thief:     a.thief,
		Card:      a.card,
		Clickable: a.clickable,
	}
	j.Row, j.Column = a.row, a.column
	return json.Marshal(j)
}

func (a *area) UnmarshalJSON(bs []byte) error {
	var j jArea
	err := json.Unmarshal(bs, &j)
	if err != nil {
		return err
	}
	a.row, a.column, a.thief, a.card, a.clickable = j.Row, j.Column, j.Thief, j.Card, j.Clickable
	return nil
}

type areaID struct {
	row    int
	column int
}

type jAreaID struct {
	Row    int `json:"row" binding:"min=1,max=8"`
	Column int `json:"column" binding:"min=1,max=8"`
}

func (aid areaID) MarshalJSON() ([]byte, error) {
	j := jAreaID{Row: aid.row, Column: aid.column}
	return json.Marshal(j)
}

func (aid *areaID) UnmarshalJSON(bs []byte) error {
	var j jAreaID
	err := json.Unmarshal(bs, &j)
	if err != nil {
		return err
	}
	aid.row, aid.column = j.Row, j.Column
	return nil
}

func newArea(row, col int, card card) area {
	return area{areaID: areaID{row: row, column: col}, card: card}
}

func (a area) hasThief() bool {
	return a.thief.pid != pidNone
}

func (a area) hasCard() bool {
	return a.card.kind != cdNone
}

func hasArea(as []area, a2 area) bool {
	for _, a1 := range as {
		b := a1.row == a2.row && a1.column == a2.column
		if b {
			return true
		}
	}
	return false
}

func (a area) hasOtherThief(p player) bool {
	return a.hasThief() && a.thief.pid != p.ID
}

type thief struct {
	pid  int
	from areaID
}

type jThief struct {
	PID  int    `json:"pid"`
	From areaID `json:"from"`
}

func (t thief) MarshalJSON() ([]byte, error) {
	j := jThief{
		PID:  t.pid,
		From: t.from,
	}
	return json.Marshal(j)
}

func (t *thief) UnmarshalJSON(bs []byte) error {
	var j jThief
	err := json.Unmarshal(bs, &j)
	if err != nil {
		return err
	}
	t.pid, t.from = j.PID, j.From
	return nil
}
