package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("card", func() {

	Describe("newCard", func() {
		var (
			cd     card
			kind   cdKind
			facing cdFacing
		)

		JustBeforeEach(func() {
			cd = newCard(kind, facing)
		})

		Context("when face up lamp", func() {

			BeforeEach(func() {
				kind, facing = cdLamp, cdFaceUp
			})

			It("should create card", func() {
				Expect(cd).ToNot(BeNil())
			})

			It("should create a card of proper type", func() {
				Expect(cd.kind).To(Equal(kind))
			})

			It("should create a card of proper facing", func() {
				Expect(cd.facing).To(Equal(facing))
			})
		})
	})

	Describe("newDeck", func() {

		var deck []card

		JustBeforeEach(func() {
			deck = newDeck()
		})

		It("should have the correct length", func() {
			Expect(deck).To(HaveLen(64))
		})

		It("should have eight lamps", func() {
			lamp := func(c card) bool { return c.kind == cdLamp }
			Expect(countBy(deck, lamp)).To(Equal(8))
		})

		It("should have eight camels", func() {
			camel := func(c card) bool { return c.kind == cdCamel }
			Expect(countBy(deck, camel)).To(Equal(8))
		})
	})

	DescribeTable("String",
		func(k cdKind, expected string) {
			Expect(k.String()).To(Equal(expected))
		},
		Entry("Lamp", cdLamp, "Lamp"),
		Entry("Camel", cdCamel, "Camel"),
		Entry("Sword", cdSword, "Sword"),
		Entry("Carpet", cdCarpet, "Carpet"),
		Entry("Coins", cdCoins, "Coins"),
		Entry("Turban", cdTurban, "Turban"),
		Entry("Jewels", cdJewels, "Jewels"),
		Entry("Guard", cdGuard, "Guard"),
		Entry("Starter Camel", cdSCamel, "Camel"),
		Entry("Starter Lamp", cdSLamp, "Lamp"),
		Entry("Default", cdNone, "None"),
	)

	DescribeTable("kindFor",
		func(s string, expected cdKind) {
			Expect(kindFor(s)).To(Equal(expected))
		},
		Entry("Lamp", "Lamp", cdLamp),
		Entry("Camel", "Camel", cdCamel),
		Entry("Sword", "Sword", cdSword),
		Entry("Carpet", "Carpet", cdCarpet),
		Entry("Coins", "Coins", cdCoins),
		Entry("Turban", "Turban", cdTurban),
		Entry("Jewels", "Jewels", cdJewels),
		Entry("Guard", "Guard", cdGuard),
		Entry("Starter Camel", "Start-Camel", cdSCamel),
		Entry("Starter Lamp", "Start-Lamp", cdSLamp),
		Entry("Default", "Nonw", cdNone),
	)

	DescribeTable("idString",
		func(k cdKind, expected string) {
			Expect(k.idString()).To(Equal(expected))
		},
		Entry("Lamp", cdLamp, "lamp"),
		Entry("Camel", cdCamel, "camel"),
		Entry("Sword", cdSword, "sword"),
		Entry("Carpet", cdCarpet, "carpet"),
		Entry("Coins", cdCoins, "coins"),
		Entry("Turban", cdTurban, "turban"),
		Entry("Jewels", cdJewels, "jewels"),
		Entry("Guard", cdGuard, "guard"),
		Entry("Starter Camel", cdSCamel, "start-camel"),
		Entry("Starter Lamp", cdSLamp, "start-lamp"),
		Entry("Default", cdNone, "none"),
	)

	type valueCase struct {
		Card     card
		Expected int
	}

	DescribeTable("value",
		func(c valueCase) {
			Expect(c.Card.value()).To(Equal(c.Expected))
		},
		Entry("Lamp", valueCase{
			Card:     newCard(cdLamp, cdFaceUp),
			Expected: 1,
		}),
		Entry("Camel", valueCase{
			Card:     newCard(cdCamel, cdFaceUp),
			Expected: 4,
		}),
		Entry("Sword", valueCase{
			Card:     newCard(cdSword, cdFaceUp),
			Expected: 5,
		}),
		Entry("Carpet", valueCase{
			Card:     newCard(cdCarpet, cdFaceUp),
			Expected: 3,
		}),
		Entry("Coins", valueCase{
			Card:     newCard(cdCoins, cdFaceUp),
			Expected: 3,
		}),
		Entry("Turban", valueCase{
			Card:     newCard(cdTurban, cdFaceUp),
			Expected: 2,
		}),
		Entry("Jewels", valueCase{
			Card:     newCard(cdJewels, cdFaceUp),
			Expected: 2,
		}),
		Entry("Guard", valueCase{
			Card:     newCard(cdGuard, cdFaceUp),
			Expected: -1,
		}),
		Entry("Starter Camel", valueCase{
			Card:     newCard(cdSCamel, cdFaceUp),
			Expected: 0,
		}),
		Entry("Starter Lamp", valueCase{
			Card:     newCard(cdSLamp, cdFaceUp),
			Expected: 0,
		}),
		Entry("Default", valueCase{
			Card:     newCard(cdNone, cdFaceUp),
			Expected: 0,
		}),
	)

	type turnCase struct {
		Cards     []card
		Direction cdFacing
		Expected  []card
	}

	DescribeTable("turn",
		func(c turnCase) {
			Expect(turn(c.Direction, c.Cards)).To(Equal(c.Expected))
		},
		Entry("Empty", turnCase{
			Cards:     noCards,
			Direction: cdFaceUp,
			Expected:  noCards,
		}),
		Entry("FaceUp One", turnCase{
			Cards:     []card{newCard(cdCamel, cdFaceDown)},
			Direction: cdFaceUp,
			Expected:  []card{newCard(cdCamel, cdFaceUp)},
		}),
		Entry("FaceUp Two", turnCase{
			Cards:     []card{newCard(cdCamel, cdFaceDown), newCard(cdLamp, cdFaceDown)},
			Direction: cdFaceUp,
			Expected:  []card{newCard(cdCamel, cdFaceUp), newCard(cdLamp, cdFaceUp)},
		}),
		Entry("FaceUp Three", turnCase{
			Cards:     []card{newCard(cdCamel, cdFaceDown), newCard(cdLamp, cdFaceUp), newCard(cdJewels, cdFaceDown)},
			Direction: cdFaceUp,
			Expected:  []card{newCard(cdCamel, cdFaceUp), newCard(cdLamp, cdFaceUp), newCard(cdJewels, cdFaceUp)},
		}),
		Entry("FaceDown Three", turnCase{
			Cards:     []card{newCard(cdCamel, cdFaceDown), newCard(cdLamp, cdFaceUp), newCard(cdJewels, cdFaceDown)},
			Direction: cdFaceDown,
			Expected:  []card{newCard(cdCamel, cdFaceDown), newCard(cdLamp, cdFaceDown), newCard(cdJewels, cdFaceDown)},
		}),
	)

	type unmarshalCase struct {
		Data     []byte
		Expected card
		Error    bool
	}

	DescribeTable("UnmarshalJSON",
		func(c unmarshalCase) {
			var cd card

			err := cd.UnmarshalJSON(c.Data)
			if c.Error {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).ToNot(HaveOccurred())
				Expect(cd).To(Equal(c.Expected))
			}
		},
		Entry("FaceUp Camel", unmarshalCase{
			Data:     []byte(`{ "kind": "camel", "facing": 1 }`),
			Expected: newCard(cdCamel, cdFaceUp),
			Error:    false,
		}),
		Entry("FaceDown Jewels", unmarshalCase{
			Data:     []byte(`{ "kind": "jewels", "facing": 0 }`),
			Expected: newCard(cdJewels, cdFaceDown),
			Error:    false,
		}),
		Entry("FaceDown Staring Lampe", unmarshalCase{
			Data:     []byte(`{ "kind": "start-lamp", "facing": 0 }`),
			Expected: newCard(cdSLamp, cdFaceDown),
			Error:    false,
		}),
		Entry("Invalid JSON", unmarshalCase{
			Data:     []byte(`{ "kind": start-lamp, "facing": 0 }`),
			Expected: newCard(cdSLamp, cdFaceDown),
			Error:    true,
		}),
	)

	type unmarshalKindCase struct {
		Data     []byte
		Expected cdKind
		Error    bool
	}

	DescribeTable("cdKind UnmarshalJSON",
		func(c unmarshalKindCase) {
			var kind cdKind

			err := kind.UnmarshalJSON(c.Data)
			if c.Error {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).ToNot(HaveOccurred())
				Expect(kind).To(Equal(c.Expected))
			}
		},
		Entry("Camel", unmarshalKindCase{
			Data:     []byte(`"camel"`),
			Expected: cdCamel,
			Error:    false,
		}),
		Entry("Starting Lamp", unmarshalKindCase{
			Data:     []byte(`"start-lamp"`),
			Expected: cdSLamp,
			Error:    false,
		}),
		Entry("Invalid JSON", unmarshalKindCase{
			Data:     []byte(`{ => start-lamp }`),
			Expected: cdSLamp,
			Error:    true,
		}),
	)

	type marshalCase struct {
		Card     card
		Expected string
	}

	DescribeTable("MarshalJSON",
		func(c marshalCase) {
			v, err := c.Card.MarshalJSON()
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(MatchJSON(c.Expected))
		},
		Entry("FaceUp Camel", marshalCase{
			Card:     newCard(cdCamel, cdFaceUp),
			Expected: `{ "kind": "camel", "facing": 1 }`,
		}),
		Entry("FaceDown Jewels", marshalCase{
			Card:     newCard(cdJewels, cdFaceDown),
			Expected: `{ "kind": "jewels", "facing": 0 }`,
		}),
		Entry("FaceDown Starting Lamp", marshalCase{
			Card:     newCard(cdSLamp, cdFaceDown),
			Expected: `{ "kind": "start-lamp", "facing": 0 }`,
		}),
	)

	type findIndexCase struct {
		cards []card
		test  func(card) bool
		index int
		found bool
	}

	DescribeTable("findIndexFor",
		func(c findIndexCase) {
			i, found := findIndexFor(c.cards, c.test)
			if c.found {
				Expect(found).To(BeTrue())
				Expect(i).To(Equal(c.index))
			} else {
				Expect(found).To(BeFalse())
			}
		},
		Entry("No Cards", findIndexCase{
			cards: noCards,
			test:  func(cd card) bool { return true },
			index: -1,
			found: false,
		}),
		Entry("Found", findIndexCase{
			cards: []card{newCard(cdLamp, cdFaceDown), newCard(cdCamel, cdFaceDown), newCard(cdJewels, cdFaceDown)},
			test:  func(cd card) bool { return cd.kind == cdCamel },
			index: 1,
			found: true,
		}),
		Entry("Not Found", findIndexCase{
			cards: []card{newCard(cdLamp, cdFaceDown), newCard(cdCamel, cdFaceDown), newCard(cdJewels, cdFaceDown)},
			test:  func(cd card) bool { return cd.kind == cdCoins },
			index: -1,
			found: false,
		}),
	)

	Describe("getIndex", func() {
		type getIndexCase struct {
			req   *http.Request
			cards []card
			index int
			err   error
		}

		var (
			ctx  *gin.Context
			resp *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			resp = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(resp)
		})

		DescribeTable("getIndexTable",
			func(c getIndexCase) {

				c.req.Header.Set("Content-Type", "application/json")
				ctx.Request = c.req

				i, err := getIndex(ctx, c.cards)
				if c.err != nil {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(c.err.Error()))
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(i).To(Equal(c.index))
				}
			},
			Entry("No Cards", getIndexCase{
				req: httptest.NewRequest(
					http.MethodPost,
					"/"+playCardPath,
					strings.NewReader(`{ "kind": "camel" }`),
				),
				cards: noCards,
				index: -1,
				err:   fmt.Errorf("card not found for %q: %w", cdCamel, errValidation),
			}),
			Entry("Found", getIndexCase{
				req: httptest.NewRequest(
					http.MethodPost,
					"/"+playCardPath,
					strings.NewReader(`{ "kind": "camel" }`),
				),
				cards: []card{newCard(cdLamp, cdFaceDown), newCard(cdCamel, cdFaceDown), newCard(cdJewels, cdFaceDown)},
				index: 1,
				err:   nil,
			}),
			Entry("Not Found", getIndexCase{
				req: httptest.NewRequest(
					http.MethodPost,
					"/"+playCardPath,
					strings.NewReader(`{ "kind": "coins" }`),
				),
				cards: []card{newCard(cdLamp, cdFaceDown), newCard(cdCamel, cdFaceDown), newCard(cdJewels, cdFaceDown)},
				index: -1,
				err:   fmt.Errorf("card not found for %q: %w", cdCoins, errValidation),
			}),
			Entry("Invalid Context", getIndexCase{
				req: httptest.NewRequest(
					http.MethodPost,
					"/"+playCardPath,
					strings.NewReader(`{ "kind" => "coins" }`),
				),
				cards: []card{newCard(cdLamp, cdFaceDown), newCard(cdCamel, cdFaceDown), newCard(cdJewels, cdFaceDown)},
				index: -1,
				err:   fmt.Errorf("unable to get card index from context: %w", errValidation),
			}),
		)
	})
})
