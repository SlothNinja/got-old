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

var _ = Describe("PlaceThief", func() {
	var (
		c          *gin.Context
		cp         player
		a          area
		resp       *httptest.ResponseRecorder
		g          game
		cu, u1, u2 user.User2
		found      bool
		err        error
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()

		g = createGame(c, u1, u2)
	})

	JustBeforeEach(func() {
		g, err = g.PlaceThief(c)
	})

	Context("when no current user", func() {
		It("should indicate there is no current user", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to find current user"))
		})

		Describe("when there is a valid request", func() {

			BeforeEach(func() {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+placeThiefPath,
					strings.NewReader(`{ "row": 1 , "column": 1 }`),
				)
				c.Request.Header.Set("Content-Type", "application/json")
			})

			It("should not place thief", func() {
				a, found = g.grid.area(1, 1)
				Expect(found).To(BeTrue())
				Expect(a.hasThief()).To(BeFalse())
			})

		})
	})

	Context("when the current user is the current player", func() {

		BeforeEach(func() {
			if g.CPUserIndices[0] == 1 {
				cu = u1
			} else {
				cu = u2
			}
			user.WithCurrent(c, cu)

			cp, found = g.currentPlayerFor(cu)
			Expect(found).To(BeTrue())
		})

		Describe("when there is a valid request", func() {

			BeforeEach(func() {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+placeThiefPath,
					strings.NewReader(`{ "row": 1 , "column": 1 }`),
				)
				c.Request.Header.Set("Content-Type", "application/json")
			})

			It("should place thief", func() {
				a, found = g.grid.area(1, 1)
				Expect(found).To(BeTrue())
				Expect(a.thief.pid).Should(Equal(cp.ID))
			})

			Context("when thief already in selected area", func() {
				BeforeEach(func() {
					a, found = g.grid.area(1, 1)
					Expect(found).To(BeTrue())

					g, a = g.placeThief(cp, a)
				})

				It("should indicate area already has thief", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when no card in selected area", func() {
				BeforeEach(func() {
					a, found = g.grid.area(1, 1)
					Expect(found).To(BeTrue())

					g, a = g.removeCardFrom(a)
				})

				It("should indicate area has no card", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("has no card"))
				})
			})

			Context("when wrong phase to place thief", func() {
				BeforeEach(func() {
					g.Phase = phasePlayCard
				})

				It("should indicate wrong phase", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("wrong phase"))
				})
			})
		})

		Describe("when there is an invalid request", func() {

			BeforeEach(func() {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+placeThiefPath,
					strings.NewReader(`{ "row": -1 , "column": 1 }`),
				)
				c.Request.Header.Set("Content-Type", "application/json")
			})

			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should not place thief", func() {
				a, found = g.grid.area(1, 1)
				Expect(found).To(BeTrue())
				Expect(a.hasThief()).To(BeFalse())
			})
		})
	})

	Describe("when the current user is not the current player", func() {

		BeforeEach(func() {
			c.Request = httptest.NewRequest(http.MethodPost, "/"+showPath+"/1", nil)
			c.Params = gin.Params{gin.Param{"hid", "1"}}

			if g.CPUserIndices[0] == 1 {
				user.WithCurrent(c, u2)
				cp, found = g.currentPlayerFor(u2)
				Expect(found).To(BeFalse())
			} else {
				user.WithCurrent(c, u1)
				cp, found = g.currentPlayerFor(u1)
				Expect(found).To(BeFalse())
			}
		})

		Describe("when there is a valid request", func() {

			BeforeEach(func() {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+placeThiefPath,
					strings.NewReader(`{ "row": 1 , "column": 1 }`),
				)
				c.Request.Header.Set("Content-Type", "application/json")
			})

			It("should indicate current player not found", func() {
				Expect(err.Error()).To(ContainSubstring("current player not found"))
			})

			It("should not place thief", func() {
				a, found = g.grid.area(1, 1)
				Expect(found).To(BeTrue())
				Expect(a.hasThief()).To(BeFalse())
			})
		})
	})
})
