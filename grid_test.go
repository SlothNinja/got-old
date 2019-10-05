package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"errors"

	"bitbucket.org/SlothNinja/user"
	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("areaID", func() {
	var (
		aid areaID
		err error
	)

	Describe("MarshalJSON", func() {
		var v []byte

		type entry struct {
			aid         areaID
			expected    string
			expectedErr error
		}

		DescribeTable("MarshalJSON",

			func(e entry) {
				aid = e.aid
				v, err = aid.MarshalJSON()
				if e.expectedErr != nil {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(v).To(MatchJSON(e.expected))
				}
			},
			Entry("areaID{0,0}", entry{
				aid:         areaID{0, 0},
				expected:    `{ "row": 0, "column": 0 }`,
				expectedErr: nil,
			}),
			Entry("areaID{2,3}", entry{
				aid:         areaID{2, 3},
				expected:    `{ "row": 2, "column": 3 }`,
				expectedErr: nil,
			}),
		)
	})

	Describe("UnmarshalJSON", func() {
		type entry struct {
			value       []byte
			expected    areaID
			expectedErr error
		}

		DescribeTable("UnmarshalJSON",
			func(e entry) {
				err = aid.UnmarshalJSON(e.value)
				if e.expectedErr != nil {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(aid).To(Equal(e.expected))
				}
			},
			Entry("areaID{0,0}", entry{
				value:       []byte(`{ "row": 0, "column": 0 }`),
				expected:    areaID{0, 0},
				expectedErr: nil,
			}),
			Entry("areaID{2,3}", entry{
				value:       []byte(`{ "row": 2, "column": 3 }`),
				expected:    areaID{2, 3},
				expectedErr: nil,
			}),
			Entry("invalid JSON", entry{
				value:       []byte(`{ "row" 2, "column": 3 }`),
				expected:    areaID{2, 3},
				expectedErr: errValidation,
			}),
		)
	})
})

var _ = Describe("area", func() {
	var (
		a   area
		err error
	)

	Describe("MarshalJSON", func() {
		var v []byte

		type entry struct {
			area        area
			expected    string
			expectedErr error
		}

		DescribeTable("with",

			func(e entry) {
				a = e.area
				v, err = a.MarshalJSON()
				if e.expectedErr != nil {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(v).To(MatchJSON(e.expected))
				}
			},
			Entry("zero value", entry{
				area: noArea,
				expected: `{ "row": 0,
					"column": 0,
					"card": { "kind": "none", "facing": 0 },
					"thief": {
						"pid": 0,
						"from": { "row": 0, "column": 0 } },
					"clickable": false }`,
				expectedErr: nil,
			}),
			Entry("non-zero value", entry{
				area: area{
					areaID:    areaID{row: 3, column: 2},
					card:      card{kind: cdCamel, facing: cdFaceUp},
					thief:     thief{pid: 2, from: areaID{row: 4, column: 2}},
					clickable: true},
				expected: `{ "row": 3,
					"column": 2,
					"card": { "kind": "camel", "facing": 1 },
					"thief": {
						"pid": 2,
						"from": { "row": 4, "column": 2 } },
					"clickable": true }`,
				expectedErr: nil,
			}),
		)
	})

	Describe("UnmarshalJSON", func() {
		type entry struct {
			value       []byte
			expected    area
			expectedErr error
		}

		DescribeTable("with",
			func(e entry) {
				err = a.UnmarshalJSON(e.value)
				if e.expectedErr != nil {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(a).To(Equal(e.expected))
				}
			},
			Entry("zero value", entry{
				value: []byte(`{ "row": 0,
					"column": 0,
					"card": { "kind": "none", "facing": 0 },
					"thief": {
						"pid": 0,
						"from": { "row": 0, "column": 0 } },
					"clickable": false }`),
				expected:    noArea,
				expectedErr: nil,
			}),
			Entry("non-zero value", entry{
				value: []byte(`{ "row": 1,
					"column": 5,
					"card": { "kind": "lamp", "facing": 1 },
					"thief": {
						"pid": 1,
						"from": { "row": 1, "column": 2 } },
					"clickable": true }`),
				expected: area{
					areaID:    areaID{row: 1, column: 5},
					card:      card{kind: cdLamp, facing: cdFaceUp},
					thief:     thief{pid: 1, from: areaID{row: 1, column: 2}},
					clickable: true},
				expectedErr: nil,
			}),
			Entry("invalid json", entry{
				value: []byte(`{ "row": 1,
					"column" 5,
					"card": { "kind": "lamp", "facing": 1 },
					"thief": {
						"pid": 1,
						"from": { "row": 1, "column": 2 } },
					"clickable": true }`),
				expected: area{
					areaID:    areaID{row: 1, column: 1},
					card:      card{kind: cdLamp, facing: cdFaceUp},
					thief:     thief{pid: 1, from: areaID{row: 1, column: 2}},
					clickable: true},
				expectedErr: errValidation,
			}),
		)
	})

	Describe("hasOtherThief", func() {
		var result bool
		type entry struct {
			area     area
			player   player
			expected bool
		}

		DescribeTable("getArea",
			func(e entry) {
				result = e.area.hasOtherThief(e.player)
				Expect(result).To(Equal(e.expected))
			},
			Entry("no thief in area", entry{
				area:     area{thief: thief{pid: noPID}},
				player:   player{id: 1},
				expected: false,
			}),
			Entry("has thief of player in area", entry{
				area:     area{thief: thief{pid: 2}},
				player:   player{id: 2},
				expected: false,
			}),
			Entry("has thief of other player in area", entry{
				area:     area{thief: thief{pid: 1}},
				player:   player{id: 2},
				expected: true,
			}),
		)
	})
})

var _ = Describe("grid", func() {
	var (
		c          *gin.Context
		resp       *httptest.ResponseRecorder
		err        error
		a          area
		aid        areaID
		g          game
		u1, u2, u3 user.User2
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

	})

	Describe("MarshalJSON", func() {
		var v []byte

		type entry struct {
			grid        grid
			expected    string
			expectedErr error
		}

		DescribeTable("MarshalJSON",

			func(e entry) {
				v, err = json.Marshal(e.grid)
				if e.expectedErr != nil {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(v).To(MatchJSON(e.expected))
				}
			},
			Entry("2x2 non-zero value grid", entry{
				grid: grid{
					[]area{
						area{
							areaID:    areaID{row: 1, column: 1},
							card:      card{kind: cdCamel, facing: cdFaceUp},
							thief:     thief{pid: 0, from: areaID{row: 0, column: 0}},
							clickable: false,
						},
						area{
							areaID:    areaID{row: 1, column: 2},
							card:      card{kind: cdCamel, facing: cdFaceUp},
							thief:     thief{pid: 0, from: areaID{row: 0, column: 0}},
							clickable: true,
						},
					},
					[]area{
						area{
							areaID:    areaID{row: 2, column: 1},
							card:      card{kind: cdLamp, facing: cdFaceUp},
							thief:     thief{pid: 0, from: areaID{row: 0, column: 0}},
							clickable: false,
						},
						area{
							areaID:    areaID{row: 2, column: 2},
							card:      card{kind: cdLamp, facing: cdFaceUp},
							thief:     thief{pid: 0, from: areaID{row: 0, column: 0}},
							clickable: false,
						},
					},
				},
				expected: `[ [ { "row": 1,
					"column": 1,
					"card": { "kind": "camel", "facing": 1 },
					"thief": { "pid": 0, "from": { "row": 0, "column": 0 } },
					"clickable": false },
					{ "row": 1,
					"column": 2,
					"card": { "kind": "camel", "facing": 1 },
					"thief": { "pid": 0, "from": { "row": 0, "column": 0 } },
					"clickable": true } ],
					[ { "row": 2,
					"column": 1,
					"card": { "kind": "lamp", "facing": 1 },
					"thief": { "pid": 0, "from": { "row": 0, "column": 0 } },
					"clickable": false },
					{ "row": 2,
					"column": 2,
					"card": { "kind": "lamp", "facing": 1 },
					"thief": { "pid": 0, "from": { "row": 0, "column": 0 } },
					"clickable": false } ] ]`,
				expectedErr: nil,
			}),
		)
	})

	Describe("UnmarshalJSON", func() {

		var grd grid

		type entry struct {
			value       []byte
			expected    grid
			expectedErr error
		}

		DescribeTable("with",
			func(e entry) {
				err = json.Unmarshal(e.value, &grd)
				if e.expectedErr != nil {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(grd).To(Equal(e.expected))
				}
			},
			Entry("2x2 non-zero value grid", entry{
				expected: grid{
					[]area{
						area{
							areaID:    areaID{row: 1, column: 1},
							card:      card{kind: cdCamel, facing: cdFaceUp},
							thief:     thief{pid: 0, from: areaID{row: 0, column: 0}},
							clickable: false,
						},
						area{
							areaID:    areaID{row: 1, column: 2},
							card:      card{kind: cdCamel, facing: cdFaceUp},
							thief:     thief{pid: 0, from: areaID{row: 0, column: 0}},
							clickable: true,
						},
					},
					[]area{
						area{
							areaID:    areaID{row: 2, column: 1},
							card:      card{kind: cdLamp, facing: cdFaceUp},
							thief:     thief{pid: 0, from: areaID{row: 0, column: 0}},
							clickable: false,
						},
						area{
							areaID:    areaID{row: 2, column: 2},
							card:      card{kind: cdLamp, facing: cdFaceUp},
							thief:     thief{pid: 0, from: areaID{row: 0, column: 0}},
							clickable: false,
						},
					},
				},
				value: []byte(`[ [ { "row": 1,
					"column": 1,
					"card": { "kind": "camel", "facing": 1 },
					"thief": { "pid": 0, "from": { "row": 0, "column": 0 } },
					"clickable": false },
					{ "row": 1,
					"column": 2,
					"card": { "kind": "camel", "facing": 1 },
					"thief": { "pid": 0, "from": { "row": 0, "column": 0 } },
					"clickable": true } ],
					[ { "row": 2,
					"column": 1,
					"card": { "kind": "lamp", "facing": 1 },
					"thief": { "pid": 0, "from": { "row": 0, "column": 0 } },
					"clickable": false },
					{ "row": 2,
					"column": 2,
					"card": { "kind": "lamp", "facing": 1 },
					"thief": { "pid": 0, "from": { "row": 0, "column": 0 } },
					"clickable": false } ] ]`),
				expectedErr: nil,
			}),
		)
	})

	Context("when two players", func() {

		BeforeEach(func() {
			u1, u2 = createUsers()
			g = createGame(c, u1, u2)
		})

		Describe("area", func() {

			var a area

			type entry struct {
				row, column int
			}

			DescribeTable("with",

				func(e entry) {
					a = g.grid.area(e.row, e.column)
					if a != noArea {
						Expect(a.row).To(Equal(e.row))
						Expect(a.column).To(Equal(e.column))
					}
				},
				Entry("area{0,0}", entry{row: 0, column: 0}),
				Entry("area{0,0}", entry{row: 0, column: 1}),
				Entry("area{0,0}", entry{row: 0, column: 2}),
				Entry("area{0,0}", entry{row: 0, column: 3}),
				Entry("area{0,0}", entry{row: 0, column: 4}),
				Entry("area{0,0}", entry{row: 0, column: 5}),
				Entry("area{0,0}", entry{row: 0, column: 6}),
				Entry("area{0,0}", entry{row: 0, column: 7}),
				Entry("area{0,0}", entry{row: 0, column: 8}),
				Entry("area{0,0}", entry{row: 0, column: 9}),

				Entry("area{1,0}", entry{row: 1, column: 0}),
				Entry("area{1,0}", entry{row: 1, column: 1}),
				Entry("area{1,0}", entry{row: 1, column: 2}),
				Entry("area{1,0}", entry{row: 1, column: 3}),
				Entry("area{1,0}", entry{row: 1, column: 4}),
				Entry("area{1,0}", entry{row: 1, column: 5}),
				Entry("area{1,0}", entry{row: 1, column: 6}),
				Entry("area{1,0}", entry{row: 1, column: 7}),
				Entry("area{1,0}", entry{row: 1, column: 8}),
				Entry("area{1,0}", entry{row: 1, column: 9}),

				Entry("area{2,0}", entry{row: 2, column: 0}),
				Entry("area{2,0}", entry{row: 2, column: 1}),
				Entry("area{2,0}", entry{row: 2, column: 2}),
				Entry("area{2,0}", entry{row: 2, column: 3}),
				Entry("area{2,0}", entry{row: 2, column: 4}),
				Entry("area{2,0}", entry{row: 2, column: 5}),
				Entry("area{2,0}", entry{row: 2, column: 6}),
				Entry("area{2,0}", entry{row: 2, column: 7}),
				Entry("area{2,0}", entry{row: 2, column: 8}),
				Entry("area{2,0}", entry{row: 2, column: 9}),

				Entry("area{3,0}", entry{row: 3, column: 0}),
				Entry("area{3,0}", entry{row: 3, column: 1}),
				Entry("area{3,0}", entry{row: 3, column: 2}),
				Entry("area{3,0}", entry{row: 3, column: 3}),
				Entry("area{3,0}", entry{row: 3, column: 4}),
				Entry("area{3,0}", entry{row: 3, column: 5}),
				Entry("area{3,0}", entry{row: 3, column: 6}),
				Entry("area{3,0}", entry{row: 3, column: 7}),
				Entry("area{3,0}", entry{row: 3, column: 8}),
				Entry("area{3,0}", entry{row: 3, column: 9}),

				Entry("area{4,0}", entry{row: 4, column: 0}),
				Entry("area{4,0}", entry{row: 4, column: 1}),
				Entry("area{4,0}", entry{row: 4, column: 2}),
				Entry("area{4,0}", entry{row: 4, column: 3}),
				Entry("area{4,0}", entry{row: 4, column: 4}),
				Entry("area{4,0}", entry{row: 4, column: 5}),
				Entry("area{4,0}", entry{row: 4, column: 6}),
				Entry("area{4,0}", entry{row: 4, column: 7}),
				Entry("area{4,0}", entry{row: 4, column: 8}),
				Entry("area{4,0}", entry{row: 4, column: 9}),

				Entry("area{5,0}", entry{row: 5, column: 0}),
				Entry("area{5,0}", entry{row: 5, column: 1}),
				Entry("area{5,0}", entry{row: 5, column: 2}),
				Entry("area{5,0}", entry{row: 5, column: 3}),
				Entry("area{5,0}", entry{row: 5, column: 4}),
				Entry("area{5,0}", entry{row: 5, column: 5}),
				Entry("area{5,0}", entry{row: 5, column: 6}),
				Entry("area{5,0}", entry{row: 5, column: 7}),
				Entry("area{5,0}", entry{row: 5, column: 8}),
				Entry("area{5,0}", entry{row: 5, column: 9}),

				Entry("area{6,0}", entry{row: 6, column: 0}),
				Entry("area{6,0}", entry{row: 6, column: 1}),
				Entry("area{6,0}", entry{row: 6, column: 2}),
				Entry("area{6,0}", entry{row: 6, column: 3}),
				Entry("area{6,0}", entry{row: 6, column: 4}),
				Entry("area{6,0}", entry{row: 6, column: 5}),
				Entry("area{6,0}", entry{row: 6, column: 6}),
				Entry("area{6,0}", entry{row: 6, column: 7}),
				Entry("area{6,0}", entry{row: 6, column: 8}),
				Entry("area{6,0}", entry{row: 6, column: 9}),

				Entry("area{7,0}", entry{row: 7, column: 0}),
				Entry("area{7,0}", entry{row: 7, column: 1}),
				Entry("area{7,0}", entry{row: 7, column: 2}),
				Entry("area{7,0}", entry{row: 7, column: 3}),
				Entry("area{7,0}", entry{row: 7, column: 4}),
				Entry("area{7,0}", entry{row: 7, column: 5}),
				Entry("area{7,0}", entry{row: 7, column: 6}),
				Entry("area{7,0}", entry{row: 7, column: 7}),
				Entry("area{7,0}", entry{row: 7, column: 8}),
				Entry("area{7,0}", entry{row: 7, column: 9}),
			)
		})

		Describe("getArea", func() {

			DescribeTable("with",

				func(row, column int, expectedAreaID areaID, expectedErr error) {
					c.Request = httptest.NewRequest(
						http.MethodPost,
						"/"+placeThiefPath,
						strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`, row, column)),
					)
					c.Request.Header.Set("Content-Type", "application/json")
					a, err = g.getArea(c)
					if expectedErr != nil {
						Expect(errors.Is(err, errValidation)).To(BeTrue())
					} else {
						Expect(err).ToNot(HaveOccurred())
						Expect(a.areaID).To(Equal(expectedAreaID))
					}
				},
				Entry("areaID{0,0}", 0, 0, areaID{0, 0}, errValidation),
				Entry("areaID{0,1}", 0, 1, areaID{0, 1}, errValidation),
				Entry("areaID{0,2}", 0, 2, areaID{0, 2}, errValidation),
				Entry("areaID{0,3}", 0, 3, areaID{0, 3}, errValidation),
				Entry("areaID{0,4}", 0, 4, areaID{0, 4}, errValidation),
				Entry("areaID{0,5}", 0, 5, areaID{0, 5}, errValidation),
				Entry("areaID{0,6}", 0, 6, areaID{0, 6}, errValidation),
				Entry("areaID{0,7}", 0, 7, areaID{0, 7}, errValidation),
				Entry("areaID{0,8}", 0, 8, areaID{0, 8}, errValidation),
				Entry("areaID{0,9}", 0, 9, areaID{0, 9}, errValidation),

				Entry("areaID{1,0}", 1, 0, areaID{1, 0}, errValidation),
				Entry("areaID{1,1}", 1, 1, areaID{1, 1}, nil),
				Entry("areaID{1,2}", 1, 2, areaID{1, 2}, nil),
				Entry("areaID{1,3}", 1, 3, areaID{1, 3}, nil),
				Entry("areaID{1,4}", 1, 4, areaID{1, 4}, nil),
				Entry("areaID{1,5}", 1, 5, areaID{1, 5}, nil),
				Entry("areaID{1,6}", 1, 6, areaID{1, 6}, nil),
				Entry("areaID{1,7}", 1, 7, areaID{1, 7}, nil),
				Entry("areaID{1,8}", 1, 8, areaID{1, 8}, nil),
				Entry("areaID{1,9}", 1, 9, areaID{1, 9}, errValidation),

				Entry("areaID{2,0}", 2, 0, areaID{2, 0}, errValidation),
				Entry("areaID{2,1}", 2, 1, areaID{2, 1}, nil),
				Entry("areaID{2,2}", 2, 2, areaID{2, 2}, nil),
				Entry("areaID{2,3}", 2, 3, areaID{2, 3}, nil),
				Entry("areaID{2,4}", 2, 4, areaID{2, 4}, nil),
				Entry("areaID{2,5}", 2, 5, areaID{2, 5}, nil),
				Entry("areaID{2,6}", 2, 6, areaID{2, 6}, nil),
				Entry("areaID{2,7}", 2, 7, areaID{2, 7}, nil),
				Entry("areaID{2,8}", 2, 8, areaID{2, 8}, nil),
				Entry("areaID{2,9}", 2, 9, areaID{2, 9}, errValidation),

				Entry("areaID{3,0}", 3, 0, areaID{3, 0}, errValidation),
				Entry("areaID{3,1}", 3, 1, areaID{3, 1}, nil),
				Entry("areaID{3,2}", 3, 2, areaID{3, 2}, nil),
				Entry("areaID{3,3}", 3, 3, areaID{3, 3}, nil),
				Entry("areaID{3,4}", 3, 4, areaID{3, 4}, nil),
				Entry("areaID{3,5}", 3, 5, areaID{3, 5}, nil),
				Entry("areaID{3,6}", 3, 6, areaID{3, 6}, nil),
				Entry("areaID{3,7}", 3, 7, areaID{3, 7}, nil),
				Entry("areaID{3,8}", 3, 8, areaID{3, 8}, nil),
				Entry("areaID{3,9}", 3, 9, areaID{3, 9}, errValidation),

				Entry("areaID{4,0}", 4, 0, areaID{4, 0}, errValidation),
				Entry("areaID{4,1}", 4, 1, areaID{4, 1}, nil),
				Entry("areaID{4,2}", 4, 2, areaID{4, 2}, nil),
				Entry("areaID{4,3}", 4, 3, areaID{4, 3}, nil),
				Entry("areaID{4,4}", 4, 4, areaID{4, 4}, nil),
				Entry("areaID{4,5}", 4, 5, areaID{4, 5}, nil),
				Entry("areaID{4,6}", 4, 6, areaID{4, 6}, nil),
				Entry("areaID{4,7}", 4, 7, areaID{4, 7}, nil),
				Entry("areaID{4,8}", 4, 8, areaID{4, 8}, nil),
				Entry("areaID{4,9}", 4, 9, areaID{4, 9}, errValidation),

				Entry("areaID{5,0}", 5, 0, areaID{5, 0}, errValidation),
				Entry("areaID{5,1}", 5, 1, areaID{5, 1}, nil),
				Entry("areaID{5,2}", 5, 2, areaID{5, 2}, nil),
				Entry("areaID{5,3}", 5, 3, areaID{5, 3}, nil),
				Entry("areaID{5,4}", 5, 4, areaID{5, 4}, nil),
				Entry("areaID{5,5}", 5, 5, areaID{5, 5}, nil),
				Entry("areaID{5,6}", 5, 6, areaID{5, 6}, nil),
				Entry("areaID{5,7}", 5, 7, areaID{5, 7}, nil),
				Entry("areaID{5,8}", 5, 8, areaID{5, 8}, nil),
				Entry("areaID{5,9}", 5, 9, areaID{5, 9}, errValidation),

				Entry("areaID{6,0}", 6, 0, areaID{6, 0}, errValidation),
				Entry("areaID{6,1}", 6, 1, areaID{6, 1}, nil),
				Entry("areaID{6,2}", 6, 2, areaID{6, 2}, nil),
				Entry("areaID{6,3}", 6, 3, areaID{6, 3}, nil),
				Entry("areaID{6,4}", 6, 4, areaID{6, 4}, nil),
				Entry("areaID{6,5}", 6, 5, areaID{6, 5}, nil),
				Entry("areaID{6,6}", 6, 6, areaID{6, 6}, nil),
				Entry("areaID{6,7}", 6, 7, areaID{6, 7}, nil),
				Entry("areaID{6,8}", 6, 8, areaID{6, 8}, nil),
				Entry("areaID{6,9}", 6, 9, areaID{6, 9}, errValidation),

				Entry("areaID{7,0}", 7, 0, areaID{7, 0}, errValidation),
				Entry("areaID{7,1}", 7, 1, areaID{7, 1}, errValidation),
				Entry("areaID{7,2}", 7, 2, areaID{7, 2}, errValidation),
				Entry("areaID{7,3}", 7, 3, areaID{7, 3}, errValidation),
				Entry("areaID{7,4}", 7, 4, areaID{7, 4}, errValidation),
				Entry("areaID{7,5}", 7, 5, areaID{7, 5}, errValidation),
				Entry("areaID{7,6}", 7, 6, areaID{7, 6}, errValidation),
				Entry("areaID{7,7}", 7, 7, areaID{7, 7}, errValidation),
				Entry("areaID{7,8}", 7, 8, areaID{7, 8}, errValidation),
				Entry("areaID{7,9}", 7, 9, areaID{7, 9}, errValidation),
			)
		})

		Describe("getAreaID", func() {
			DescribeTable("with",

				func(row, column int, expectedAreaID areaID, expectedErr error) {
					c.Request = httptest.NewRequest(
						http.MethodPost,
						"/"+placeThiefPath,
						strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`, row, column)),
					)
					c.Request.Header.Set("Content-Type", "application/json")
					aid, err = g.getAreaID(c)
					if expectedErr != nil {
						Expect(errors.Is(err, errValidation)).To(BeTrue())
					} else {
						Expect(err).ToNot(HaveOccurred())
						Expect(aid).To(Equal(expectedAreaID))
					}
				},
				Entry("areaID{0,0}", 0, 0, areaID{0, 0}, errValidation),
				Entry("areaID{0,1}", 0, 1, areaID{0, 1}, errValidation),
				Entry("areaID{0,2}", 0, 2, areaID{0, 2}, errValidation),
				Entry("areaID{0,3}", 0, 3, areaID{0, 3}, errValidation),
				Entry("areaID{0,4}", 0, 4, areaID{0, 4}, errValidation),
				Entry("areaID{0,5}", 0, 5, areaID{0, 5}, errValidation),
				Entry("areaID{0,6}", 0, 6, areaID{0, 6}, errValidation),
				Entry("areaID{0,7}", 0, 7, areaID{0, 7}, errValidation),
				Entry("areaID{0,8}", 0, 8, areaID{0, 8}, errValidation),
				Entry("areaID{0,9}", 0, 9, areaID{0, 9}, errValidation),

				Entry("areaID{1,0}", 1, 0, areaID{1, 0}, errValidation),
				Entry("areaID{1,1}", 1, 1, areaID{1, 1}, nil),
				Entry("areaID{1,2}", 1, 2, areaID{1, 2}, nil),
				Entry("areaID{1,3}", 1, 3, areaID{1, 3}, nil),
				Entry("areaID{1,4}", 1, 4, areaID{1, 4}, nil),
				Entry("areaID{1,5}", 1, 5, areaID{1, 5}, nil),
				Entry("areaID{1,6}", 1, 6, areaID{1, 6}, nil),
				Entry("areaID{1,7}", 1, 7, areaID{1, 7}, nil),
				Entry("areaID{1,8}", 1, 8, areaID{1, 8}, nil),
				Entry("areaID{1,9}", 1, 9, areaID{1, 9}, errValidation),

				Entry("areaID{2,0}", 2, 0, areaID{2, 0}, errValidation),
				Entry("areaID{2,1}", 2, 1, areaID{2, 1}, nil),
				Entry("areaID{2,2}", 2, 2, areaID{2, 2}, nil),
				Entry("areaID{2,3}", 2, 3, areaID{2, 3}, nil),
				Entry("areaID{2,4}", 2, 4, areaID{2, 4}, nil),
				Entry("areaID{2,5}", 2, 5, areaID{2, 5}, nil),
				Entry("areaID{2,6}", 2, 6, areaID{2, 6}, nil),
				Entry("areaID{2,7}", 2, 7, areaID{2, 7}, nil),
				Entry("areaID{2,8}", 2, 8, areaID{2, 8}, nil),
				Entry("areaID{2,9}", 2, 9, areaID{2, 9}, errValidation),

				Entry("areaID{3,0}", 3, 0, areaID{3, 0}, errValidation),
				Entry("areaID{3,1}", 3, 1, areaID{3, 1}, nil),
				Entry("areaID{3,2}", 3, 2, areaID{3, 2}, nil),
				Entry("areaID{3,3}", 3, 3, areaID{3, 3}, nil),
				Entry("areaID{3,4}", 3, 4, areaID{3, 4}, nil),
				Entry("areaID{3,5}", 3, 5, areaID{3, 5}, nil),
				Entry("areaID{3,6}", 3, 6, areaID{3, 6}, nil),
				Entry("areaID{3,7}", 3, 7, areaID{3, 7}, nil),
				Entry("areaID{3,8}", 3, 8, areaID{3, 8}, nil),
				Entry("areaID{3,9}", 3, 9, areaID{3, 9}, errValidation),

				Entry("areaID{4,0}", 4, 0, areaID{4, 0}, errValidation),
				Entry("areaID{4,1}", 4, 1, areaID{4, 1}, nil),
				Entry("areaID{4,2}", 4, 2, areaID{4, 2}, nil),
				Entry("areaID{4,3}", 4, 3, areaID{4, 3}, nil),
				Entry("areaID{4,4}", 4, 4, areaID{4, 4}, nil),
				Entry("areaID{4,5}", 4, 5, areaID{4, 5}, nil),
				Entry("areaID{4,6}", 4, 6, areaID{4, 6}, nil),
				Entry("areaID{4,7}", 4, 7, areaID{4, 7}, nil),
				Entry("areaID{4,8}", 4, 8, areaID{4, 8}, nil),
				Entry("areaID{4,9}", 4, 9, areaID{4, 9}, errValidation),

				Entry("areaID{5,0}", 5, 0, areaID{5, 0}, errValidation),
				Entry("areaID{5,1}", 5, 1, areaID{5, 1}, nil),
				Entry("areaID{5,2}", 5, 2, areaID{5, 2}, nil),
				Entry("areaID{5,3}", 5, 3, areaID{5, 3}, nil),
				Entry("areaID{5,4}", 5, 4, areaID{5, 4}, nil),
				Entry("areaID{5,5}", 5, 5, areaID{5, 5}, nil),
				Entry("areaID{5,6}", 5, 6, areaID{5, 6}, nil),
				Entry("areaID{5,7}", 5, 7, areaID{5, 7}, nil),
				Entry("areaID{5,8}", 5, 8, areaID{5, 8}, nil),
				Entry("areaID{5,9}", 5, 9, areaID{5, 9}, errValidation),

				Entry("areaID{6,0}", 6, 0, areaID{6, 0}, errValidation),
				Entry("areaID{6,1}", 6, 1, areaID{6, 1}, nil),
				Entry("areaID{6,2}", 6, 2, areaID{6, 2}, nil),
				Entry("areaID{6,3}", 6, 3, areaID{6, 3}, nil),
				Entry("areaID{6,4}", 6, 4, areaID{6, 4}, nil),
				Entry("areaID{6,5}", 6, 5, areaID{6, 5}, nil),
				Entry("areaID{6,6}", 6, 6, areaID{6, 6}, nil),
				Entry("areaID{6,7}", 6, 7, areaID{6, 7}, nil),
				Entry("areaID{6,8}", 6, 8, areaID{6, 8}, nil),
				Entry("areaID{6,9}", 6, 9, areaID{6, 9}, errValidation),

				Entry("areaID{7,0}", 7, 0, areaID{7, 0}, errValidation),
				Entry("areaID{7,1}", 7, 1, areaID{7, 1}, errValidation),
				Entry("areaID{7,2}", 7, 2, areaID{7, 2}, errValidation),
				Entry("areaID{7,3}", 7, 3, areaID{7, 3}, errValidation),
				Entry("areaID{7,4}", 7, 4, areaID{7, 4}, errValidation),
				Entry("areaID{7,5}", 7, 5, areaID{7, 5}, errValidation),
				Entry("areaID{7,6}", 7, 6, areaID{7, 6}, errValidation),
				Entry("areaID{7,7}", 7, 7, areaID{7, 7}, errValidation),
				Entry("areaID{7,8}", 7, 8, areaID{7, 8}, errValidation),
				Entry("areaID{7,9}", 7, 9, areaID{7, 9}, errValidation),
			)
		})
	})

	Describe("numCols", func() {
		Context("with initialized grid", func() {
			It("should return number columns", func() {
				Expect(g.grid.numCols()).To(Equal(col8))
			})
		})

		Context("with empty grid", func() {
			It("should return zero", func() {
				Expect(grid{}.numCols()).To(BeZero())
			})
		})
	})

	Context("when three players", func() {

		BeforeEach(func() {
			u1, u2, u3 = createUsers3()
			g = createGame3(c, u1, u2, u3)
		})

		Describe("area", func() {

			var a area

			type entry struct {
				row, column int
			}

			DescribeTable("with",

				func(e entry) {
					a = g.grid.area(e.row, e.column)
					if a != noArea {
						Expect(a.row).To(Equal(e.row))
						Expect(a.column).To(Equal(e.column))
					}
				},
				Entry("area{0,0}", entry{row: 0, column: 0}),
				Entry("area{0,0}", entry{row: 0, column: 1}),
				Entry("area{0,0}", entry{row: 0, column: 2}),
				Entry("area{0,0}", entry{row: 0, column: 3}),
				Entry("area{0,0}", entry{row: 0, column: 4}),
				Entry("area{0,0}", entry{row: 0, column: 5}),
				Entry("area{0,0}", entry{row: 0, column: 6}),
				Entry("area{0,0}", entry{row: 0, column: 7}),
				Entry("area{0,0}", entry{row: 0, column: 8}),
				Entry("area{0,0}", entry{row: 0, column: 9}),

				Entry("area{1,0}", entry{row: 1, column: 0}),
				Entry("area{1,0}", entry{row: 1, column: 1}),
				Entry("area{1,0}", entry{row: 1, column: 2}),
				Entry("area{1,0}", entry{row: 1, column: 3}),
				Entry("area{1,0}", entry{row: 1, column: 4}),
				Entry("area{1,0}", entry{row: 1, column: 5}),
				Entry("area{1,0}", entry{row: 1, column: 6}),
				Entry("area{1,0}", entry{row: 1, column: 7}),
				Entry("area{1,0}", entry{row: 1, column: 8}),
				Entry("area{1,0}", entry{row: 1, column: 9}),

				Entry("area{2,0}", entry{row: 2, column: 0}),
				Entry("area{2,0}", entry{row: 2, column: 1}),
				Entry("area{2,0}", entry{row: 2, column: 2}),
				Entry("area{2,0}", entry{row: 2, column: 3}),
				Entry("area{2,0}", entry{row: 2, column: 4}),
				Entry("area{2,0}", entry{row: 2, column: 5}),
				Entry("area{2,0}", entry{row: 2, column: 6}),
				Entry("area{2,0}", entry{row: 2, column: 7}),
				Entry("area{2,0}", entry{row: 2, column: 8}),
				Entry("area{2,0}", entry{row: 2, column: 9}),

				Entry("area{3,0}", entry{row: 3, column: 0}),
				Entry("area{3,0}", entry{row: 3, column: 1}),
				Entry("area{3,0}", entry{row: 3, column: 2}),
				Entry("area{3,0}", entry{row: 3, column: 3}),
				Entry("area{3,0}", entry{row: 3, column: 4}),
				Entry("area{3,0}", entry{row: 3, column: 5}),
				Entry("area{3,0}", entry{row: 3, column: 6}),
				Entry("area{3,0}", entry{row: 3, column: 7}),
				Entry("area{3,0}", entry{row: 3, column: 8}),
				Entry("area{3,0}", entry{row: 3, column: 9}),

				Entry("area{4,0}", entry{row: 4, column: 0}),
				Entry("area{4,0}", entry{row: 4, column: 1}),
				Entry("area{4,0}", entry{row: 4, column: 2}),
				Entry("area{4,0}", entry{row: 4, column: 3}),
				Entry("area{4,0}", entry{row: 4, column: 4}),
				Entry("area{4,0}", entry{row: 4, column: 5}),
				Entry("area{4,0}", entry{row: 4, column: 6}),
				Entry("area{4,0}", entry{row: 4, column: 7}),
				Entry("area{4,0}", entry{row: 4, column: 8}),
				Entry("area{4,0}", entry{row: 4, column: 9}),

				Entry("area{5,0}", entry{row: 5, column: 0}),
				Entry("area{5,0}", entry{row: 5, column: 1}),
				Entry("area{5,0}", entry{row: 5, column: 2}),
				Entry("area{5,0}", entry{row: 5, column: 3}),
				Entry("area{5,0}", entry{row: 5, column: 4}),
				Entry("area{5,0}", entry{row: 5, column: 5}),
				Entry("area{5,0}", entry{row: 5, column: 6}),
				Entry("area{5,0}", entry{row: 5, column: 7}),
				Entry("area{5,0}", entry{row: 5, column: 8}),
				Entry("area{5,0}", entry{row: 5, column: 9}),

				Entry("area{6,0}", entry{row: 6, column: 0}),
				Entry("area{6,0}", entry{row: 6, column: 1}),
				Entry("area{6,0}", entry{row: 6, column: 2}),
				Entry("area{6,0}", entry{row: 6, column: 3}),
				Entry("area{6,0}", entry{row: 6, column: 4}),
				Entry("area{6,0}", entry{row: 6, column: 5}),
				Entry("area{6,0}", entry{row: 6, column: 6}),
				Entry("area{6,0}", entry{row: 6, column: 7}),
				Entry("area{6,0}", entry{row: 6, column: 8}),
				Entry("area{6,0}", entry{row: 6, column: 9}),

				Entry("area{7,0}", entry{row: 7, column: 0}),
				Entry("area{7,0}", entry{row: 7, column: 1}),
				Entry("area{7,0}", entry{row: 7, column: 2}),
				Entry("area{7,0}", entry{row: 7, column: 3}),
				Entry("area{7,0}", entry{row: 7, column: 4}),
				Entry("area{7,0}", entry{row: 7, column: 5}),
				Entry("area{7,0}", entry{row: 7, column: 6}),
				Entry("area{7,0}", entry{row: 7, column: 7}),
				Entry("area{7,0}", entry{row: 7, column: 8}),
				Entry("area{7,0}", entry{row: 7, column: 9}),

				Entry("area{8,0}", entry{row: 8, column: 0}),
				Entry("area{8,0}", entry{row: 8, column: 1}),
				Entry("area{8,0}", entry{row: 8, column: 2}),
				Entry("area{8,0}", entry{row: 8, column: 3}),
				Entry("area{8,0}", entry{row: 8, column: 4}),
				Entry("area{8,0}", entry{row: 8, column: 5}),
				Entry("area{8,0}", entry{row: 8, column: 6}),
				Entry("area{8,0}", entry{row: 8, column: 7}),
				Entry("area{8,0}", entry{row: 8, column: 8}),
				Entry("area{8,0}", entry{row: 8, column: 9}),
			)
		})

		DescribeTable("getAreaID",

			func(row, column int, expectedAreaID areaID, expectedErr error) {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+placeThiefPath,
					strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`, row, column)),
				)
				c.Request.Header.Set("Content-Type", "application/json")
				aid, err = g.getAreaID(c)
				if expectedErr != nil {
					Expect(errors.Is(err, errValidation)).To(BeTrue())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(aid).To(Equal(expectedAreaID))
				}
			},
			Entry("areaID{0,0}", 0, 0, areaID{0, 0}, errValidation),
			Entry("areaID{0,1}", 0, 1, areaID{0, 1}, errValidation),
			Entry("areaID{0,2}", 0, 2, areaID{0, 2}, errValidation),
			Entry("areaID{0,3}", 0, 3, areaID{0, 3}, errValidation),
			Entry("areaID{0,4}", 0, 4, areaID{0, 4}, errValidation),
			Entry("areaID{0,5}", 0, 5, areaID{0, 5}, errValidation),
			Entry("areaID{0,6}", 0, 6, areaID{0, 6}, errValidation),
			Entry("areaID{0,7}", 0, 7, areaID{0, 7}, errValidation),
			Entry("areaID{0,8}", 0, 8, areaID{0, 8}, errValidation),
			Entry("areaID{0,9}", 0, 9, areaID{0, 9}, errValidation),

			Entry("areaID{1,0}", 1, 0, areaID{1, 0}, errValidation),
			Entry("areaID{1,1}", 1, 1, areaID{1, 1}, nil),
			Entry("areaID{1,2}", 1, 2, areaID{1, 2}, nil),
			Entry("areaID{1,3}", 1, 3, areaID{1, 3}, nil),
			Entry("areaID{1,4}", 1, 4, areaID{1, 4}, nil),
			Entry("areaID{1,5}", 1, 5, areaID{1, 5}, nil),
			Entry("areaID{1,6}", 1, 6, areaID{1, 6}, nil),
			Entry("areaID{1,7}", 1, 7, areaID{1, 7}, nil),
			Entry("areaID{1,8}", 1, 8, areaID{1, 8}, nil),
			Entry("areaID{1,9}", 1, 9, areaID{1, 9}, errValidation),

			Entry("areaID{2,0}", 2, 0, areaID{2, 0}, errValidation),
			Entry("areaID{2,1}", 2, 1, areaID{2, 1}, nil),
			Entry("areaID{2,2}", 2, 2, areaID{2, 2}, nil),
			Entry("areaID{2,3}", 2, 3, areaID{2, 3}, nil),
			Entry("areaID{2,4}", 2, 4, areaID{2, 4}, nil),
			Entry("areaID{2,5}", 2, 5, areaID{2, 5}, nil),
			Entry("areaID{2,6}", 2, 6, areaID{2, 6}, nil),
			Entry("areaID{2,7}", 2, 7, areaID{2, 7}, nil),
			Entry("areaID{2,8}", 2, 8, areaID{2, 8}, nil),
			Entry("areaID{2,9}", 2, 9, areaID{2, 9}, errValidation),

			Entry("areaID{3,0}", 3, 0, areaID{3, 0}, errValidation),
			Entry("areaID{3,1}", 3, 1, areaID{3, 1}, nil),
			Entry("areaID{3,2}", 3, 2, areaID{3, 2}, nil),
			Entry("areaID{3,3}", 3, 3, areaID{3, 3}, nil),
			Entry("areaID{3,4}", 3, 4, areaID{3, 4}, nil),
			Entry("areaID{3,5}", 3, 5, areaID{3, 5}, nil),
			Entry("areaID{3,6}", 3, 6, areaID{3, 6}, nil),
			Entry("areaID{3,7}", 3, 7, areaID{3, 7}, nil),
			Entry("areaID{3,8}", 3, 8, areaID{3, 8}, nil),
			Entry("areaID{3,9}", 3, 9, areaID{3, 9}, errValidation),

			Entry("areaID{4,0}", 4, 0, areaID{4, 0}, errValidation),
			Entry("areaID{4,1}", 4, 1, areaID{4, 1}, nil),
			Entry("areaID{4,2}", 4, 2, areaID{4, 2}, nil),
			Entry("areaID{4,3}", 4, 3, areaID{4, 3}, nil),
			Entry("areaID{4,4}", 4, 4, areaID{4, 4}, nil),
			Entry("areaID{4,5}", 4, 5, areaID{4, 5}, nil),
			Entry("areaID{4,6}", 4, 6, areaID{4, 6}, nil),
			Entry("areaID{4,7}", 4, 7, areaID{4, 7}, nil),
			Entry("areaID{4,8}", 4, 8, areaID{4, 8}, nil),
			Entry("areaID{4,9}", 4, 9, areaID{4, 9}, errValidation),

			Entry("areaID{5,0}", 5, 0, areaID{5, 0}, errValidation),
			Entry("areaID{5,1}", 5, 1, areaID{5, 1}, nil),
			Entry("areaID{5,2}", 5, 2, areaID{5, 2}, nil),
			Entry("areaID{5,3}", 5, 3, areaID{5, 3}, nil),
			Entry("areaID{5,4}", 5, 4, areaID{5, 4}, nil),
			Entry("areaID{5,5}", 5, 5, areaID{5, 5}, nil),
			Entry("areaID{5,6}", 5, 6, areaID{5, 6}, nil),
			Entry("areaID{5,7}", 5, 7, areaID{5, 7}, nil),
			Entry("areaID{5,8}", 5, 8, areaID{5, 8}, nil),
			Entry("areaID{5,9}", 5, 9, areaID{5, 9}, errValidation),

			Entry("areaID{6,0}", 6, 0, areaID{6, 0}, errValidation),
			Entry("areaID{6,1}", 6, 1, areaID{6, 1}, nil),
			Entry("areaID{6,2}", 6, 2, areaID{6, 2}, nil),
			Entry("areaID{6,3}", 6, 3, areaID{6, 3}, nil),
			Entry("areaID{6,4}", 6, 4, areaID{6, 4}, nil),
			Entry("areaID{6,5}", 6, 5, areaID{6, 5}, nil),
			Entry("areaID{6,6}", 6, 6, areaID{6, 6}, nil),
			Entry("areaID{6,7}", 6, 7, areaID{6, 7}, nil),
			Entry("areaID{6,8}", 6, 8, areaID{6, 8}, nil),
			Entry("areaID{6,9}", 6, 9, areaID{6, 9}, errValidation),

			Entry("areaID{7,0}", 7, 0, areaID{7, 0}, errValidation),
			Entry("areaID{7,1}", 7, 1, areaID{7, 1}, nil),
			Entry("areaID{7,2}", 7, 2, areaID{7, 2}, nil),
			Entry("areaID{7,3}", 7, 3, areaID{7, 3}, nil),
			Entry("areaID{7,4}", 7, 4, areaID{7, 4}, nil),
			Entry("areaID{7,5}", 7, 5, areaID{7, 5}, nil),
			Entry("areaID{7,6}", 7, 6, areaID{7, 6}, nil),
			Entry("areaID{7,7}", 7, 7, areaID{7, 7}, nil),
			Entry("areaID{7,8}", 7, 8, areaID{7, 8}, nil),
			Entry("areaID{7,9}", 7, 9, areaID{7, 9}, errValidation),

			Entry("areaID{8,0}", 8, 0, areaID{8, 0}, errValidation),
			Entry("areaID{8,1}", 8, 1, areaID{8, 1}, errValidation),
			Entry("areaID{8,2}", 8, 2, areaID{8, 2}, errValidation),
			Entry("areaID{8,3}", 8, 3, areaID{8, 3}, errValidation),
			Entry("areaID{8,4}", 8, 4, areaID{8, 4}, errValidation),
			Entry("areaID{8,5}", 8, 5, areaID{8, 5}, errValidation),
			Entry("areaID{8,6}", 8, 6, areaID{8, 6}, errValidation),
			Entry("areaID{8,7}", 8, 7, areaID{8, 7}, errValidation),
			Entry("areaID{8,8}", 8, 8, areaID{8, 8}, errValidation),
			Entry("areaID{8,9}", 8, 9, areaID{8, 9}, errValidation),
		)
	})
})

var _ = Describe("thief MarshalJSON", func() {
	var (
		t   thief
		v   []byte
		err error
	)

	type case1 struct {
		pid         int
		from        areaID
		expected    string
		expectedErr error
	}

	DescribeTable("MarshalJSON",

		func(e case1) {
			t = thief{pid: e.pid, from: e.from}
			v, err = t.MarshalJSON()
			if e.expectedErr != nil {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).ToNot(HaveOccurred())
				Expect(v).To(MatchJSON(e.expected))
			}
		},
		Entry("thief{0, areaID{0,0}}", case1{
			pid:         0,
			from:        areaID{0, 0},
			expected:    `{ "pid": 0, "from": { "row": 0, "column": 0 } }`,
			expectedErr: nil,
		}),
		Entry("thief{1, areaID{2,3}}", case1{
			pid:         1,
			from:        areaID{2, 3},
			expected:    `{ "pid": 1, "from": { "row": 2, "column": 3 } }`,
			expectedErr: nil,
		}),
		Entry("thief{1, areaID{2,3}}", case1{
			pid:         1,
			from:        areaID{2, 3},
			expected:    `{ "pid": 1, "from": { "row": 2, "column": 3 } }`,
			expectedErr: nil,
		}),
	)
})

var _ = Describe("thief UnmarshalJSON", func() {
	var (
		t   thief
		err error
	)

	type entry struct {
		value       []byte
		expected    thief
		expectedErr error
	}

	DescribeTable("Unmarshal",
		func(e entry) {
			err = t.UnmarshalJSON(e.value)
			if e.expectedErr != nil {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).ToNot(HaveOccurred())
				Expect(t).To(Equal(e.expected))
			}
		},
		Entry("thief{0, areaID{0,0}}", entry{
			value:       []byte(`{ "pid": 0, "from": { "row": 0, "column": 0 } }`),
			expected:    thief{0, areaID{0, 0}},
			expectedErr: nil,
		}),
		Entry("thief{1, areaID{2,3}}", entry{
			value:       []byte(`{ "pid": 1, "from": { "row": 2, "column": 3 } }`),
			expected:    thief{1, areaID{2, 3}},
			expectedErr: nil,
		}),
		Entry("invalid JSON", entry{
			value:       []byte(`{ "pid" => 1, "from": { "row": 2, "column": 3 } }`),
			expected:    thief{1, areaID{2, 3}},
			expectedErr: errValidation,
		}),
	)
})
