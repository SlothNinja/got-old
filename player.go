package main

import (
	"encoding/json"
	"math/rand"

	"bitbucket.org/SlothNinja/color"
	"bitbucket.org/SlothNinja/user"
	"cloud.google.com/go/datastore"
)

var defaultColors = []color.Color{color.Yellow, color.Purple, color.Green, color.Black}

const noPID = 0

var noPlayer = player{}

// Player represents one of the players of the game.
type player struct {
	id              int
	performedAction bool
	score           int
	passed          bool
	colors          []color.Color
	user            user.User2
	hand            []card
	drawPile        []card
	discardPile     []card
}

type jPlayer struct {
	ID              int           `json:"id"`
	PerformedAction bool          `json:"performedAction"`
	Score           int           `json:"score"`
	Passed          bool          `json:"passed"`
	Colors          []color.Color `json:"colors"`
	User            user.User2    `json:"user"`
	Hand            []card        `json:"hand"`
	DrawPile        []card        `json:"draw"`
	DiscardPile     []card        `json:"discard"`
}

func (p player) MarshalJSON() ([]byte, error) {
	j := jPlayer{
		ID:              p.id,
		PerformedAction: p.performedAction,
		Score:           p.score,
		Passed:          p.passed,
		Colors:          p.colors,
		User:            p.user,
		Hand:            p.hand,
		DrawPile:        p.drawPile,
		DiscardPile:     p.discardPile,
	}
	return json.Marshal(j)
}

func (p *player) UnmarshalJSON(bs []byte) error {
	var j jPlayer
	err := json.Unmarshal(bs, &j)
	if err != nil {
		return err
	}
	p.id, p.performedAction, p.score, p.passed = j.ID, j.PerformedAction, j.Score, j.Passed
	p.colors, p.user, p.hand, p.drawPile, p.discardPile = j.Colors, j.User, j.Hand, j.DrawPile, j.DiscardPile
	return nil
}

func newPlayer() player {
	return player{
		hand:        startHand(),
		drawPile:    make([]card, 0),
		discardPile: make([]card, 0),
	}
}

func (p player) defeated(p2 player) (bool, bool) {
	if p.score != p2.score {
		return p.score > p2.score, false
	}

	lampTest := func(c card) bool { return c.kind == cdLamp }
	lamp1, lamp2 := countBy(p.hand, lampTest), countBy(p2.hand, lampTest)
	if lamp1 != lamp2 {
		return lamp1 > lamp2, false
	}

	camelTest := func(c card) bool { return c.kind == cdCamel }
	camel1, camel2 := countBy(p.hand, camelTest), countBy(p2.hand, camelTest)
	if camel1 != camel2 {
		return camel1 > camel2, false
	}

	cds1, cds2 := len(p.hand), len(p2.hand)
	if cds1 != cds2 {
		return cds1 > cds2, false
	}
	return false, true
}

// func (p player) beginningOfTurnReset() {
// 	p.clearActions()
// }

func (p player) clearActions() player {
	p.performedAction = false
	return p
}

func (p player) draw() (player, card, bool) {
	shuffle := len(p.drawPile) < 1
	if shuffle {
		p = p.shuffle()
	}

	var cd card
	p.drawPile, cd = draw(p.drawPile)
	cd.turn(cdFaceUp)
	p.hand = append(p.hand, cd)
	return p, cd, shuffle
}

func (p player) shuffle() player {
	p.drawPile, p.discardPile = p.discardPile, make([]card, 0)
	p.drawPile = turn(cdFaceDown, p.drawPile)
	rand.Shuffle(len(p.drawPile), func(i, j int) {
		p.drawPile[i], p.drawPile[j] = p.drawPile[j], p.drawPile[i]
	})
	return p
}

// Equal assume two players with the same ID are equal.
func (p player) Equal(p2 player) bool {
	return p.id == p2.id
}

func (p player) collectCards() {
	p.hand = append(p.hand, p.discardPile...)
	p.hand = append(p.hand, p.drawPile...)
	turn(cdFaceUp, p.hand)
	p.discardPile, p.drawPile = make([]card, 0), make([]card, 0)
	return
}

// Players is a slice of players of the game.
// type Players []Player

// // ToProperty implements interface for serializing players to datastore.
// func (ps Players) ToProperty() (datastore.Property, error) {
// 	v, err := json.Marshal(ps)
// 	return datastore.MkPropertyNI(string(v)), err
// }
//
// // FromProperty implements interface for serializing players from datastore.
// func (ps Players) FromProperty(prop datastore.Property) error {
// 	return json.Unmarshal([]byte(prop.Value().(string)), ps)
// }

// func UserDataFor(ps Players) (ud []user.Data) {
// 	ud = make([]user.Data, len(ps))
// 	for i := range ps {
// 		ud[i] = ps[i].User.Data
// 	}
// 	return
// }

func allPassed(ps []player) bool {
	for _, p := range ps {
		if !p.passed {
			return false
		}
	}
	return true
}

func pids(ps []player) []int {
	pids := make([]int, len(ps))
	for i := range ps {
		pids[i] = ps[i].id
	}
	return pids
}

// func uids(ps []Player) []string {
// 	uids := make([]string, len(ps))
// 	for i := range ps {
// 		uids[i] = ps[i].User.ID()
// 	}
// 	return uids
// }

func playerUKeys(ps []player) []*datastore.Key {
	ks := make([]*datastore.Key, len(ps))
	for i := range ps {
		ks[i] = ps[i].user.Key
	}
	return ks
}

func playerFindIndex(p player, ps []player) (int, bool) {
	for index := range ps {
		if p.id == ps[index].id {
			return index, true
		}
	}
	return -1, false
}

func playerByID(pid int, ps []player) player {
	for _, p := range ps {
		if p.id == pid {
			return p
		}
	}
	return noPlayer
}

// shuffle randomizes the order of players.
// note, using custom shuffle until AppEngine upgrades to 1.10
// At which time, the below shuffle can be replaced with rand.shuffle
func playerShuffle(ps []player) {
	rand.Shuffle(len(ps), func(i, j int) {
		ps[i], ps[j] = ps[j], ps[i]
	})
}

// // ByIndex returns the player at the index i in the ring of players ps
// // Wraps-around based on numPlayers.
func playerByIndex(numPlayers int, ps []player) player {
	l := len(ps)
	r := numPlayers % l
	if r < 0 {
		return ps[l+r]
	}
	return ps[r]
}

// func PlayerHideCards(ps []Player) []Player {
// 	ps2 := make([]Player, len(ps))
// 	for i, p := range ps {
// 		p2 := p.hideCards()
// 		ps2[i] = &(p2)
// 	}
// 	return ps2
// }

func (p player) hideCards() player {
	p.hand, p.discardPile, p.drawPile = noCards, noCards, noCards
	return p
}

func (g game) updatePlayer(p player) game {
	i, found := playerFindIndex(p, g.players)
	if found {
		g.players[i] = p
	}
	return g
}
