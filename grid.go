package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"bitbucket.org/SlothNinja/log"
	"github.com/gin-gonic/gin"
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

func lastRowFor(numPlayers int) int {
	if numPlayers == 2 {
		return rowF
	}
	return rowG
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

func (g grid) area(row, col int) area {
	if row < rowA || col < col1 || row > g.numRows() || col > g.numCols() {
		return noArea
	}
	return g[row-1][col-1]
}

func (g grid) numRows() int {
	return len(g)
}

func (g grid) numCols() int {
	if g.numRows() > 1 {
		return len(g[0])
	}
	return 0
}

func (g grid) each(f func(a area) area) {
	for row := range g {
		for col := range g[row] {
			g[row][col] = f(g[row][col])
		}
	}
}

func (g grid) updateArea(a area) grid {
	if a.row < rowA || a.column < col1 || a.row > g.numRows() || a.column > g.numCols() {
		log.Warningf("area: %#v not in grid", a)
		return g
	}
	g[a.row-1][a.column-1] = a
	return g
}

type area struct {
	areaID
	thief     thief
	card      card
	clickable bool
}

var noArea = area{}

func (a area) empty() area {
	return area{areaID: a.areaID}
}

func (a area) isEmpty() bool {
	return !a.hasThief() && !a.hasCard()
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

var noAreaID = areaID{}
var noAreaIDS = []areaID{}

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
	return a.thief.pid != noPID
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
	log.Warningf("a.hasThief: %v", a.hasThief())
	log.Warningf("a.thief.pid: %v", a.thief.pid)
	log.Warningf("p.id: %v", p.id)
	log.Warningf("p: %#v", p)
	log.Warningf("a: %#v", a)
	return a.hasThief() && a.thief.pid != p.id
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

func (g game) getArea(c *gin.Context) (area, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	aid, err := g.getAreaID(c)
	if err != nil {
		return noArea, err
	}

	return g.grid.area(aid.row, aid.column), nil
}

func (g game) getAreaID(c *gin.Context) (areaID, error) {
	if g.NumPlayers == 2 {
		obj := struct {
			Row    int `json:"row" binding:"min=1,max=6"`
			Column int `json:"column" binding:"min=1,max=8"`
		}{}

		err := c.ShouldBindJSON(&obj)
		if err != nil {
			return noAreaID, fmt.Errorf("%s: %w", err, errValidation)
		}
		return areaID{row: obj.Row, column: obj.Column}, nil
	}

	obj := struct {
		Row    int `json:"row" binding:"min=1,max=7"`
		Column int `json:"column" binding:"min=1,max=8"`
	}{}

	err := c.ShouldBindJSON(&obj)
	if err != nil {
		return noAreaID, fmt.Errorf("%s: %w", err, errValidation)
	}

	return areaID{row: obj.Row, column: obj.Column}, nil
}
