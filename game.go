// Package got implements the card game, Guild of Thieves.
package main

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/SlothNinja/color"
	"bitbucket.org/SlothNinja/gtype"
	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/stack"
	"bitbucket.org/SlothNinja/user"
	"cloud.google.com/go/datastore"
)

const (
	kind     = "Game"
	msgExit  = "Exiting"
	msgEnter = "Entering"
	gameKey  = "game"
)

type game struct {
	Key        *datastore.Key `datastore:"-"`
	Encoded    string         `datastore:",noindex"`
	EncodedLog string         `datastore:",noindex"`
	Log        `datastore:"-"`
	state      `datastore:"-"`
	stack.Stack
	Header
}

var noGame = game{}

type state struct {
	players        []player
	grid           grid
	jewels         card
	stepped        int
	playedCard     card
	selectedAreaID areaID
}

type jState struct {
	Players        []player `json:"players"`
	Grid           grid     `json:"grid"`
	Jewels         card     `json:"jewels"`
	Stepped        int      `json:"stepped"`
	PlayedCard     card     `json:"playedCard"`
	SelectedAreaID areaID   `json:"selectedAreaID"`
}

func (s state) MarshalJSON() ([]byte, error) {
	j := jState{
		Players:        s.players,
		Grid:           s.grid,
		Jewels:         s.jewels,
		Stepped:        s.stepped,
		PlayedCard:     s.playedCard,
		SelectedAreaID: s.selectedAreaID,
	}
	return json.Marshal(j)
}

func (s *state) UnmarshalJSON(bs []byte) error {
	var j jState
	err := json.Unmarshal(bs, &j)
	if err != nil {
		return err
	}
	s.players, s.grid, s.jewels, s.stepped, s.playedCard, s.selectedAreaID =
		j.Players, j.Grid, j.Jewels, j.Stepped, j.PlayedCard, j.SelectedAreaID
	return nil
}

func newGame(id int64) game {
	var g game
	g.Type = gtype.GOT
	g.Key = newKey(id)
	return g
}

func newKey(id int64) *datastore.Key {
	return datastore.IDKey(kind, id, nil)
}

func (g game) ID() int64 {
	return g.Key.ID
}

type jGame struct {
	ID     int64          `json:"id"`
	Key    *datastore.Key `json:"key"`
	Log    Log            `json:"log"`
	Stack  stack.Stack    `json:"undoStack"`
	State  state          `json:"state"`
	Header Header         `json:"header"`
}

func (g game) MarshalJSON() ([]byte, error) {
	return json.Marshal(jGame{
		ID:     g.ID(),
		Key:    g.Key,
		Log:    g.Log,
		Stack:  g.Stack,
		State:  g.state,
		Header: g.Header,
	})
}

func (g *game) Load(ps []datastore.Property) error {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	err := datastore.LoadStruct(g, ps)
	if err != nil {
		return err
	}

	var s state
	err = json.Unmarshal([]byte(g.Encoded), &s)
	if err != nil {
		return err
	}
	g.state = s

	var gl Log
	err = json.Unmarshal([]byte(g.EncodedLog), &gl)
	if err != nil {
		return err
	}
	g.Log = gl
	return nil
}

func (g *game) Save() ([]datastore.Property, error) {
	encoded, err := json.Marshal(g.state)
	if err != nil {
		return nil, err
	}

	encodedLog, err := json.Marshal(g.Log)
	if err != nil {
		return nil, err
	}

	g.Encoded = string(encoded)
	g.EncodedLog = string(encodedLog)
	return datastore.SaveStruct(g)
}

func (g *game) LoadKey(k *datastore.Key) error {
	g.Key = k
	return nil
}

func (g game) withoutHistory() ([]*datastore.Key, []interface{}) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	h := newHeaderEntity(g)
	h.CreatedAt, h.UpdatedAt = g.CreatedAt, g.UpdatedAt

	return []*datastore.Key{g.Key, h.Key}, []interface{}{&g, &h}
}

func (g game) withHistory() ([]*datastore.Key, []interface{}) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	h := newHeaderEntity(g)
	h.CreatedAt, h.UpdatedAt = g.CreatedAt, g.UpdatedAt

	history := newHistory(g)
	h.CreatedAt, h.UpdatedAt = g.CreatedAt, g.UpdatedAt

	return []*datastore.Key{g.Key, h.Key, history.Key}, []interface{}{&g, &h, &history}
}

type history struct{ game }

func newHistory(g game) history {
	var h history
	h.game = g
	h.Key = newHistoryKey(g.ID(), g.Stack.Current)
	return h
}

func (h *history) Load(ps []datastore.Property) error {
	return h.game.Load(ps)
}

func (h *history) Save() ([]datastore.Property, error) {
	return h.game.Save()
}

func (h *history) LoadKey(k *datastore.Key) error {
	h.Key = k
	return nil
}

func newHistoryKey(id, count int64) *datastore.Key {
	return datastore.NameKey(kind, keyName(id, count), newKey(id))
}

func keyName(id int64, count int64) string {
	return fmt.Sprintf("%d-%d", id, count)
}

func (h history) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID     string `json:"id"`
		Header Header `json:"header"`
		State  state  `json:"state"`
		Log    Log    `json:"log"`
	}{
		ID:     h.Key.Name,
		Header: h.Header,
		State:  h.state,
		Log:    h.Log,
	})
}

func (g game) addNewPlayers() game {
	for i := 0; i < g.NumPlayers; i++ {
		g = g.addNewPlayer(i + 1)
	}
	return g
}

// CurrentPlayer returns the player whose turn it is.
func (g game) currentPlayerFor(u user.User2) player {
	if len(g.CPUserIndices) < 1 {
		return noPlayer
	}

	pid1 := g.CPUserIndices[0]
	if u.Admin {
		return playerByID(pid1, g.players)
	}

	pid2, found := g.PlayerIDFor(u.ID())
	if found && (pid1 == pid2) {
		return playerByID(pid1, g.players)
	}
	return noPlayer
}

// // SelectedThiefArea returns the area corresponding to a previously selected thief.
// func (g game) SelectedThiefArea() (area, bool) {
// 	return g.grid.area(g.selectedAreaID.row, g.selectedAreaID.column)
// }

func (g game) addNewPlayer(pid int) game {
	p := g.createPlayer(pid)
	g.players = append(g.players, p)
	return g
}

func (g game) createPlayer(pid int) player {
	p := newPlayer()
	p.id = pid
	p.user = g.Users[pid-1]

	p.colors = make([]color.Color, g.NumPlayers)

	for i := 0; i < g.NumPlayers; i++ {
		index := (i - p.id) % g.NumPlayers
		if index < 0 {
			index += g.NumPlayers
		}
		color := defaultColors[index]
		p.colors[i] = color
	}

	return p
}

func (g game) beginningOfPhaseReset() game {
	for i, p := range g.players {
		p = p.clearActions()
		p.passed = false
		g.players[i] = p
	}
	return g
}

func (g game) beginningOfTurnReset(p player) game {
	p = p.clearActions()
	p.passed = false
	return g.updatePlayer(p)
}

func (g game) endOfTurnUpdateFor(p player) game {
	g = g.updateJewels()
	g.stepped = 0
	p.hand = turn(cdFaceUp, p.hand)
	return g.updatePlayer(p)
}

func (g game) updateJewels() game {
	switch g.jewels = g.playedCard; g.jewels.kind {
	case cdSLamp:
		g.jewels = card{kind: cdLamp, facing: cdFaceUp}
	case cdSCamel:
		g.jewels = card{kind: cdCamel, facing: cdFaceUp}
	}
	return g
}
