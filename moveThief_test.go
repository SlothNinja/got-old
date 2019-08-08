package main

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"bitbucket.org/SlothNinja/user"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Move Thief", func() {
	var (
		c          *gin.Context
		cp         player
		resp       *httptest.ResponseRecorder
		g          game
		cu, u1, u2 user.User2
		found      bool
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
			cp, found = g.currentPlayerFor(cu)
			Expect(found).To(BeTrue())

			from, found = g.grid.area(from.row, from.column)
			Expect(found).To(BeTrue())
			Expect(from.thief.pid).To(BeZero())

			to, found = g.grid.area(to.row, to.column)
			Expect(found).To(BeTrue())
			Expect(to.thief.pid).To(Equal(cp.ID))
		})
	}

	JustBeforeEach(func() {
		g, err = g.MoveThief(c)
	})

	Context("when there is no current user", func() {
		BeforeEach(func() {
			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/"+placeThiefPath+"/1",
				strings.NewReader(`{ "row": 2, "column": 3 }`),
			)

			from, found = g.grid.area(1, 1)
			Expect(found).To(BeTrue())

			from.thief.pid = 1
			g = g.updateArea(from)

			g.selectedAreaID = from.areaID

			to, found = g.grid.area(2, 3)
			Expect(found).To(BeTrue())

		})

		It("should indicate there is no current user", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to find current user"))
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

			cp, found = g.currentPlayerFor(cu)
			Expect(found).To(BeTrue())

			from, found = g.grid.area(1, 1)
			Expect(found).To(BeTrue())

			from.thief.pid = cp.ID
			g = g.updateArea(from)

			g.playedCard = newCard(cdCamel, cdFaceUp)

			g.selectedAreaID = from.areaID

			to, found = g.grid.area(2, 3)
			Expect(found).To(BeTrue())

		})

		Context("when there is a valid request", func() {
			BeforeEach(func() {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+placeThiefPath,
					strings.NewReader(`{ "row": 2, "column": 3 }`),
				)
				c.Request.Header.Set("Content-Type", "application/json")
			})

			AssertSuccessfulBehavior()

		})
	})
})
