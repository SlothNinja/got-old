package main

import (
	"math/rand"

	snp "bitbucket.org/SlothNinja/player"
	"cloud.google.com/go/datastore"
)

const pidNone = snp.NoPID

// Player represents one of the players of the game.
type player struct {
	snp.Player
	Hand        []card `json:"hand"`
	DrawPile    []card `json:"draw"`
	DiscardPile []card `json:"discard"`
}

// type playerNoCards struct {
// 	player.Player
// 	User        *user.User2 `json:"user"`
// 	Hand        []card     `json:"-"`
// 	DrawPile    []card     `json:"-"`
// 	DiscardPile []card     `json:"-"`
// }

func newPlayer() player {
	return player{
		Hand:        startHand(),
		DrawPile:    make([]card, 0),
		DiscardPile: make([]card, 0),
	}
}

func (p player) defeated(p2 player) (bool, bool) {
	if p.Score != p2.Score {
		return p.Score > p2.Score, false
	}

	lampTest := func(c card) bool { return c.kind == cdLamp }
	lamp1, lamp2 := countBy(p.Hand, lampTest), countBy(p2.Hand, lampTest)
	if lamp1 != lamp2 {
		return lamp1 > lamp2, false
	}

	camelTest := func(c card) bool { return c.kind == cdCamel }
	camel1, camel2 := countBy(p.Hand, camelTest), countBy(p2.Hand, camelTest)
	if camel1 != camel2 {
		return camel1 > camel2, false
	}

	cds1, cds2 := len(p.Hand), len(p2.Hand)
	if cds1 != cds2 {
		return cds1 > cds2, false
	}
	return false, true
}

// func (p player) beginningOfTurnReset() {
// 	p.clearActions()
// }

func (p player) clearActions() player {
	p.PerformedAction = false
	return p
}

func (p player) draw() (player, card, bool) {
	shuffle := len(p.DrawPile) < 1
	if shuffle {
		p = p.shuffle()
	}

	var cd card
	p.DrawPile, cd = draw(p.DrawPile)
	cd.turn(cdFaceUp)
	p.Hand = append(p.Hand, cd)
	return p, cd, shuffle
}

func (p player) shuffle() player {
	p.DrawPile, p.DiscardPile = p.DiscardPile, make([]card, 0)
	p.DrawPile = turn(cdFaceDown, p.DrawPile)
	rand.Shuffle(len(p.DrawPile), func(i, j int) {
		p.DrawPile[i], p.DrawPile[j] = p.DrawPile[j], p.DrawPile[i]
	})
	return p
}

// Equal assume two players with the same ID are equal.
func (p player) Equal(p2 player) bool {
	return p.ID == p2.ID
}

func (p player) collectCards() {
	p.Hand = append(p.Hand, p.DiscardPile...)
	p.Hand = append(p.Hand, p.DrawPile...)
	turn(cdFaceUp, p.Hand)
	p.DiscardPile, p.DrawPile = make([]card, 0), make([]card, 0)
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
		if !p.Passed {
			return false
		}
	}
	return true
}

func pids(ps []player) []int {
	pids := make([]int, len(ps))
	for i := range ps {
		pids[i] = ps[i].ID
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
		ks[i] = ps[i].User.Key
	}
	return ks
}

func playerFindIndex(p player, ps []player) (int, bool) {
	for index := range ps {
		if p.ID == ps[index].ID {
			return index, true
		}
	}
	return -1, false
}

func playerByID(pid int, ps []player) (player, bool) {
	for _, p := range ps {
		if p.ID == pid {
			return p, true
		}
	}
	return player{}, false
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
	p.Hand, p.DiscardPile, p.DrawPile = []card{}, []card{}, []card{}
	return p
}

func (g game) updatePlayer(p player) game {
	i, found := playerFindIndex(p, g.players)
	if found {
		g.players[i] = p
	}
	return g
}
