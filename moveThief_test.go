package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn"
	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("s.placeThief", func() {
	var (
		c      *gin.Context
		s      server
		r      *gin.Engine
		resp   *httptest.ResponseRecorder
		req    *http.Request
		u1, u2 user.User
	)

	BeforeEach(func() {
		setGinMode()
		r = newRouter(newCookieStore())

		u1, u2 = createUsers()
		es := make(map[*datastore.Key]interface{})
		es[newKey(1)] = createGame(c, u1, u2)
		s = server{&sn.Mock{Entities: es}}
		addRoutes(rootPath, r, s)

		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		req = httptest.NewRequest(
			http.MethodPut,
			"/"+moveThiefPath+"/1",
			strings.NewReader(`{ "row": 1 , "column": 1 }`),
		)
	})

	JustBeforeEach(func() {
		r.ServeHTTP(resp, req)
	})

	Context("when no current user", func() {
		It("should indicate there is no current user", func() {
			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Body.String()).To(ContainSubstring("only the current player can perform the selected action"))
		})
	})
})

var _ = Describe("g.moveThief", func() {
	var (
		c          *gin.Context
		cp         player
		resp       *httptest.ResponseRecorder
		g          game
		cu, u1, u2 user.User
		from, to   area
		err        error
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()

		g = createGame(c, u1, u2)
		g.Phase = phaseMoveThief
	})

	AssertFailedBehavior := func() {
		It("should return an OK status", func() {
			Expect(resp.Result().StatusCode).To(Equal(http.StatusOK))
		})

		It("should provide a message", func() {
			Expect(err).ToNot(Equal(""))
		})

		It("should remain in move thief phase", func() {
			Expect(g.Phase).To(Equal(phaseMoveThief))
		})

		It("should not move thief", func() {
			Expect(from.thief.pid).ToNot(BeZero())
			Expect(to.thief.pid).To(BeZero())
		})
	}

	AssertSuccessfulBehavior := func() {
		It("should not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should proceed to claim item phase", func() {
			Expect(g.Phase).To(Equal(phaseClaimItem))
		})

		It("should move thief", func() {
			cp = g.currentPlayerFor(cu)

			from = g.grid.area(from.row, from.column)
			Expect(from.thief.pid).To(BeZero())

			to = g.grid.area(to.row, to.column)
			Expect(to.thief.pid).To(Equal(cp.id))
		})
	}

	JustBeforeEach(func() {
		g, err = g.moveThiefAction(c)
	})

	Context("when there is no current user", func() {
		BeforeEach(func() {
			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/"+placeThiefPath+"/1",
				strings.NewReader(`{ "row": 2, "column": 3 }`),
			)

			from = g.grid.area(1, 1)

			from.thief.pid = 1
			g.grid = g.grid.updateArea(from)

			g.selectedAreaID = from.areaID

			to = g.grid.area(2, 3)

		})

		It("should indicate there is no current user", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("only the current player can perform the selected action"))
		})

		AssertFailedBehavior()
	})

	Context("when the current user is the current player", func() {

		BeforeEach(func() {
			c.Request = httptest.NewRequest(http.MethodPost, "/"+moveThiefPath+"/1", nil)
			c.Params = gin.Params{gin.Param{"hid", "1"}}

			if g.CPUserIndices[0] == 1 {
				cu = u1
			} else {
				cu = u2
			}
			user.WithCurrent(c, cu)

			cp = g.currentPlayerFor(cu)

			from = g.grid.area(4, 4)

			from.thief.pid = cp.id
			g.grid = g.grid.updateArea(from)

			g.selectedAreaID = from.areaID

			to = g.grid.area(2, 3)

		})

		Context("when there is a valid request", func() {
			Context("when the played card is a camel", func() {
				BeforeEach(func() {
					to = g.grid.area(2, 3)
					g.playedCard = newCard(cdCamel, cdFaceUp)
					c.Request = httptest.NewRequest(
						http.MethodPost,
						"/"+moveThiefPath+"/1",
						strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`,
							to.row, to.column)),
					)
					c.Request.Header.Set("Content-Type", "application/json")
				})
				AssertSuccessfulBehavior()
			})

			Context("when the played card is a lamp", func() {
				BeforeEach(func() {
					to = g.grid.area(1, 4)
					g.playedCard = newCard(cdLamp, cdFaceUp)
					c.Request = httptest.NewRequest(
						http.MethodPost,
						"/"+moveThiefPath+"/1",
						strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`,
							to.row, to.column)),
					)
					c.Request.Header.Set("Content-Type", "application/json")
				})

				AssertSuccessfulBehavior()

				Context("when not first turn", func() {
					var discardPile int

					BeforeEach(func() {
						g.Turn = 10
						cp.drawPile = append(cp.drawPile, card{cdCamel, cdFaceDown}, card{cdLamp, cdFaceDown})
						discardPile = len(cp.discardPile)
						g.updatePlayer(cp)
					})

					AssertSuccessfulBehavior()

					It("should place claimed card in discard pile", func() {
						cp = playerByID(cp.id, g.players)
						Expect(cp.discardPile).To(HaveLen(discardPile + 1))
					})
				})

				Context("when no selected thief", func() {
					BeforeEach(func() {
						g.selectedAreaID = noAreaID
					})

					AssertFailedBehavior()
				})

				Context("when invalid to area", func() {
					BeforeEach(func() {
						to = g.grid.area(0, 0)
						c.Request = httptest.NewRequest(
							http.MethodPost,
							"/"+moveThiefPath+"/1",
							strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`,
								to.row, to.column)),
						)
						c.Request.Header.Set("Content-Type", "application/json")
					})

					AssertFailedBehavior()
				})

				Context("when from area has thief of another player", func() {
					BeforeEach(func() {
						from.thief.pid = cp.id + 1
						g.grid = g.grid.updateArea(from)
					})

					AssertFailedBehavior()
				})

				Context("when no card played", func() {
					BeforeEach(func() {
						g.playedCard = noCard
					})

					AssertFailedBehavior()
				})

				Context("when card does not permit movement to selected destination", func() {
					BeforeEach(func() {
						to = g.grid.area(2, 4)
						c.Request = httptest.NewRequest(
							http.MethodPost,
							"/"+moveThiefPath+"/1",
							strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`,
								to.row, to.column)),
						)
						c.Request.Header.Set("Content-Type", "application/json")
					})

					AssertFailedBehavior()
				})
			})

			Context("when the played card is a sword", func() {
				BeforeEach(func() {
					to = g.grid.area(2, 4)
					to.thief.pid = cp.id + 1
					g.grid = g.grid.updateArea(to)
					g.playedCard = newCard(cdSword, cdFaceUp)
					c.Request = httptest.NewRequest(
						http.MethodPost,
						"/"+moveThiefPath+"/1",
						strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`,
							to.row, to.column)),
					)
					c.Request.Header.Set("Content-Type", "application/json")
				})

				AssertSuccessfulBehavior()

				Context("when not first turn", func() {
					var discardPile int

					BeforeEach(func() {
						g.Turn = 10
						cp.drawPile = append(cp.drawPile, card{cdCamel, cdFaceDown}, card{cdLamp, cdFaceDown})
						discardPile = len(cp.discardPile)
						g.updatePlayer(cp)
					})

					AssertSuccessfulBehavior()

					It("should place claimed card in discard pile", func() {
						cp = playerByID(cp.id, g.players)
						Expect(cp.discardPile).To(HaveLen(discardPile + 1))
					})
				})
			})

			Context("when the played card is a coin", func() {
				BeforeEach(func() {
					to = g.grid.area(5, 4)
					cp.drawPile = append(cp.drawPile, card{cdCamel, cdFaceDown}, card{cdLamp, cdFaceDown})
					g.updatePlayer(cp)
					g.playedCard = newCard(cdCoins, cdFaceUp)
					c.Request = httptest.NewRequest(
						http.MethodPost,
						"/"+moveThiefPath+"/1",
						strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`,
							to.row, to.column)),
					)
					c.Request.Header.Set("Content-Type", "application/json")
				})
				AssertSuccessfulBehavior()
			})

			Context("when the played card is a turban", func() {
				BeforeEach(func() {
					to = g.grid.area(5, 4)
					cp.drawPile = append(cp.drawPile, card{cdCamel, cdFaceDown}, card{cdLamp, cdFaceDown})
					g.updatePlayer(cp)
					g.playedCard = newCard(cdTurban, cdFaceUp)
					c.Request = httptest.NewRequest(
						http.MethodPost,
						"/"+moveThiefPath+"/1",
						strings.NewReader(fmt.Sprintf(`{ "row": %d, "column": %d }`,
							to.row, to.column)),
					)
					c.Request.Header.Set("Content-Type", "application/json")
				})

				Context("when step yet taken", func() {
					BeforeEach(func() {
						g.stepped = 0
					})

					It("should not error", func() {
						Expect(err).ToNot(HaveOccurred())
					})

					It("should move thief", func() {
						cp = g.currentPlayerFor(cu)

						from = g.grid.area(from.row, from.column)
						Expect(from.thief.pid).To(BeZero())

						to = g.grid.area(to.row, to.column)
						Expect(to.thief.pid).To(Equal(cp.id))
					})
				})

				Context("when step taken", func() {
					BeforeEach(func() {
						g.stepped = 1
					})

					AssertSuccessfulBehavior()
				})
			})
		})
	})
})

var _ = Describe("g.bumpedTo", func() {
	var (
		c           *gin.Context
		resp        *httptest.ResponseRecorder
		g           game
		u1, u2      user.User
		a, from, to area
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()
		g = createGame(c, u1, u2)
	})

	JustBeforeEach(func() {
		a = g.grid.bumpedTo(from, to)
	})

	Context("when moving from above", func() {
		BeforeEach(func() {
			from = g.grid.area(1, 4)
			to = g.grid.area(3, 4)
		})

		It("should return space below", func() {
			Expect(a.areaID).To(Equal(areaID{to.row + 1, to.column}))
		})
	})

	Context("when moving from below", func() {
		BeforeEach(func() {
			from = g.grid.area(5, 4)
			to = g.grid.area(3, 4)
		})

		It("should return space above", func() {
			Expect(a.areaID).To(Equal(areaID{to.row - 1, to.column}))
		})
	})

	Context("when moving from the left", func() {
		BeforeEach(func() {
			from = g.grid.area(3, 1)
			to = g.grid.area(3, 4)
		})

		It("should return space to the right", func() {
			Expect(a.areaID).To(Equal(areaID{to.row, to.column + 1}))
		})
	})

	Context("when moving from the right", func() {
		BeforeEach(func() {
			from = g.grid.area(3, 6)
			to = g.grid.area(3, 4)
		})

		It("should return space to the left", func() {
			Expect(a.areaID).To(Equal(areaID{to.row, to.column - 1}))
		})
	})

	Context("default return for invalid bump", func() {
		BeforeEach(func() {
			from = g.grid.area(3, 4)
			to = g.grid.area(3, 4)
		})

		It("should return space to the left", func() {
			Expect(a.areaID).To(BeZero())
		})
	})
})
