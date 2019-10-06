package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SlothNinja/log"
	"github.com/gin-gonic/gin"
)

// card is a playing card used to form grid, player's hand, and player's deck.
type card struct {
	kind   cdKind
	facing cdFacing
}

var (
	noCard  = card{}
	noCards = []card{}
)

type jCard struct {
	Kind   cdKind   `json:"kind"`
	Facing cdFacing `json:"facing"`
}

func (c *card) UnmarshalJSON(bs []byte) error {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	var j jCard
	err := json.Unmarshal(bs, &j)
	if err != nil {
		return err
	}
	c.kind, c.facing = j.Kind, j.Facing
	return nil
}

func (c card) MarshalJSON() ([]byte, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	j := jCard{Kind: c.kind, Facing: c.facing}
	return json.Marshal(j)
}

// newCard provides a new card having the specified kind and facing.
func newCard(k cdKind, f cdFacing) card {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return card{kind: k, facing: f}
}

func newDeck() []card {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	deck := make([]card, 64)
	for j := 0; j < 8; j++ {
		for i, k := range []cdKind{cdLamp, cdCamel, cdSword, cdCarpet, cdCoins, cdTurban, cdJewels, cdGuard} {
			deck[i+j*8] = newCard(k, cdFaceDown)
		}
	}
	return deck
}

// turn sets the facing of the card to the specified value.
func (c card) turn(f cdFacing) card {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	c.facing = f
	return c
}

// turn sets the facing of the cards to the specified value.
func turn(f cdFacing, cs []card) []card {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cs2 := make([]card, len(cs))
	for i := range cs {
		cs2[i] = cs[i].turn(f)
	}
	return cs2
}

// removeAt returns a slice with the card at the specified index removed.
// removeAt will panic if index exceeds bounds of card slice.
func removeAt(i int, cs []card) []card {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return append(cs[:i], cs[i+1:]...)
}

// drawFrom returns a slice with the card removed as well as the card removed.
// drawFrom will panic if index exceeds bounds of card slice.
func drawFrom(i int, cs []card) ([]card, card) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return removeAt(i, cs), cs[i]
}

// draw retruns a slice with the first card removed as well as the card removed.
// draw will panic if index exceeds bounds of card slice.
func draw(cs []card) ([]card, card) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return drawFrom(0, cs)
}

// findIndexFor return index for the first card in the card slice that satisfies test.
// findIndexFor also returns whether a card in the card slice satisfied test.
func findIndexFor(cs []card, test func(card) bool) (int, bool) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	for i := range cs {
		found := test(cs[i])
		if found {
			return i, true
		}
	}
	return -1, false
}

// startHand returns a new starting hand of cards.
func startHand() []card {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return []card{
		newCard(cdSLamp, cdFaceUp),
		newCard(cdSLamp, cdFaceUp),
		newCard(cdSCamel, cdFaceUp),
	}
}

// value provides the point value of a card.
func (c card) value() int {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	switch c.kind {
	case cdSword:
		return 5
	case cdCamel:
		return 4
	case cdCarpet, cdCoins:
		return 3
	case cdTurban, cdJewels:
		return 2
	case cdLamp:
		return 1
	case cdGuard:
		return -1
	default:
		return 0
	}
}

// countBy returns the count of cards in the slice of cards that satisfy test.
func countBy(cs []card, test func(card) bool) int {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	count := 0
	for _, c := range cs {
		if test(c) {
			count++
		}
	}
	return count
}

// getIndex returns the index of the card in the card slice having the kind specified by "kind"
// in the JSON object received via the gin.Context.
func getIndex(c *gin.Context, cs []card) (int, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	obj := new(card)
	err := c.ShouldBindJSON(obj)
	if err != nil {
		return -1, fmt.Errorf("unable to get card index from context: %w", errValidation)
	}

	index, found := findIndexFor(cs, func(c card) bool { return c.kind == obj.kind })
	if !found {
		return -1, fmt.Errorf("card not found for %q: %w", obj.kind, errValidation)
	}
	return index, nil
}

// Kind is used to specify what kind a card has.
type cdKind int

const (
	// None indicates the card has no kind or is not present.
	cdNone cdKind = iota
	// Lamp indicates the card is a lamp card.
	cdLamp
	// Camel indicates the card is a carmel card.
	cdCamel
	// Sword indicates the card is a sword card.
	cdSword
	// Carpet indicates the card is a carpet card.
	cdCarpet
	// Coins indicates the card is a coins card.
	cdCoins
	// Turban indicates the card is a turban card.
	cdTurban
	// Jewels indicates the card is a jewels card.
	cdJewels
	// Guard indicates the card is a guard card.
	cdGuard
	// SCamel indicates the card is a starting camel card.
	cdSCamel
	// SLamp indicates the card is a starting lamp card.
	cdSLamp
)

// MarshalJSON implements the json.Marshaler interface for Kind.
func (k cdKind) MarshalJSON() ([]byte, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return json.Marshal(k.idString())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Kind.
func (k *cdKind) UnmarshalJSON(bs []byte) error {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	var s string

	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	*k = kindFor(s)
	return nil
}

// kindFor returns the Kind represented by the string.
func kindFor(s string) cdKind {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	switch s = strings.ToLower(s); s {
	case "lamp":
		return cdLamp
	case "camel":
		return cdCamel
	case "sword":
		return cdSword
	case "carpet":
		return cdCarpet
	case "coins":
		return cdCoins
	case "turban":
		return cdTurban
	case "jewels":
		return cdJewels
	case "guard":
		return cdGuard
	case "start-camel":
		return cdSCamel
	case "start-lamp":
		return cdSLamp
	default:
		return cdNone
	}
}

// String returns a string representation of the Kind.
// String does not distinguish between Camel and SCamel kinds.
// String also does not distinguish between Lamp and SLamp kinds.
func (k cdKind) String() string {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	switch k {
	case cdLamp:
		return "Lamp"
	case cdCamel:
		return "Camel"
	case cdSword:
		return "Sword"
	case cdCarpet:
		return "Carpet"
	case cdCoins:
		return "Coins"
	case cdTurban:
		return "Turban"
	case cdJewels:
		return "Jewels"
	case cdGuard:
		return "Guard"
	case cdSCamel:
		return "Camel"
	case cdSLamp:
		return "Lamp"
	default:
		return "None"
	}
}

// lString returns a lower case representation of the Kind.
func (k cdKind) lString() string {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	return strings.ToLower(k.String())
}

// idString returns a lower case representation of the Kind.
// idString distinguishes between Camel and SCamel kinds.
// idString also distinguish between Lamp and SLamp kinds.
func (k cdKind) idString() string {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	switch k {
	case cdSCamel:
		return "start-camel"
	case cdSLamp:
		return "start-lamp"
	default:
		return k.lString()
	}
}

// cdFacing is used to indicate whether card is facing up or down.
type cdFacing int

const (
	// cdFaceDown indicates card is facing down
	// cdFaceDown is also the Zero value for cdFacing.
	cdFaceDown cdFacing = iota
	// cdFaceUp indicates a card is face up.
	cdFaceUp
)
