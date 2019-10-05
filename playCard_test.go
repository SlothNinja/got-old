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

var _ = Describe("Play Card", func() {
	var (
		c             *gin.Context
		cp            player
		hand, discard int
		resp          *httptest.ResponseRecorder
		g             game
		cu, u1, u2    user.User2
		err           error
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()

		g = createGame(c, u1, u2)
		g.Phase = phasePlayCard
	})

	JustBeforeEach(func() {
		g, err = g.PlayCard(c)
	})

	AssertCUFailedBehavior := func() {
		It("should return an OK status", func() {
			Expect(resp.Result().StatusCode).To(Equal(http.StatusOK))
		})

		It("should provide a message", func() {
			Expect(err).ToNot(Equal(""))
		})

		It("should remain in play card phase", func() {
			Expect(g.Phase).To(Equal(phasePlayCard))
		})

		It("should not play card", func() {
			cp = g.currentPlayerFor(cu)

			Expect(cp.hand).To(HaveLen(hand))
			Expect(cp.discardPile).To(HaveLen(discard))
		})
	}

	AssertNoCUFailedBehavior := func() {
		It("should return an OK status", func() {
			Expect(resp.Result().StatusCode).To(Equal(http.StatusOK))
		})

		It("should provide a message", func() {
			Expect(err).ToNot(Equal(""))
		})

		It("should remain in play card phase", func() {
			Expect(g.Phase).To(Equal(phasePlayCard))
		})

	}

	AssertSuccessfulBehavior := func() {
		It("should not provide a message", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should proceed to select thief phase", func() {
			Expect(g.Phase).To(Equal(phaseSelectThief))
		})

		It("should play card", func() {
			cp = g.currentPlayerFor(cu)
			Expect(cp.hand).To(HaveLen(hand - 1))
			Expect(cp.discardPile).To(HaveLen(discard + 1))
		})
	}

	Context("when there is no current user", func() {
		BeforeEach(func() {
			c.Request = httptest.NewRequest(http.MethodPost, "/"+showPath+"/1", nil)
			c.Params = gin.Params{gin.Param{"hid", "1"}}
		})

		It("should indicate there is no current user", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("only the current player can perform the selected action"))
		})

		AssertNoCUFailedBehavior()
	})

	Context("when the current user is the current player", func() {

		BeforeEach(func() {
			c.Request = httptest.NewRequest(http.MethodPost, "/"+showPath+"/1", nil)
			c.Params = gin.Params{gin.Param{"hid", "1"}}

			if g.CPUserIndices[0] == 1 {
				cu = u1
			} else {
				cu = u2
			}
			user.WithCurrent(c, cu)

			cp = g.currentPlayerFor(cu)
			hand, discard = len(cp.hand), len(cp.discardPile)
		})

		Context("when there is a valid request", func() {

			BeforeEach(func() {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+placeThiefPath,
					strings.NewReader(`{ "kind": "start-camel" }`),
				)
				c.Request.Header.Set("Content-Type", "application/json")
			})

			AssertSuccessfulBehavior()

			Context("when wrong phase to play card", func() {
				BeforeEach(func() {
					g.Phase = phasePlaceThieves
				})

				It("should indicate wrong phase to play card", func() {
					Expect(err.Error()).To(ContainSubstring("wrong phase"))
				})
			})
		})

		Context("when there is an invalid request", func() {

			BeforeEach(func() {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+placeThiefPath,
					strings.NewReader(`{ "kind" => "start-camel" }`),
				)
				c.Request.Header.Set("Content-Type", "application/json")
			})

			AssertCUFailedBehavior()
		})

	})
})
