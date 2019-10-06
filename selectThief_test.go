package main

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Select Thief", func() {
	var (
		c          *gin.Context
		a          area
		cp         player
		resp       *httptest.ResponseRecorder
		g          game
		cu, u1, u2 user.User
		err        error
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()

		g = createGame(c, u1, u2)
	})

	AssertFailedBehavior := func() {
		It("should provide a message", func() {
			Expect(err).To(HaveOccurred())
		})

		It("should not select area", func() {
			Expect(g.selectedAreaID).To(BeZero())
		})
	}

	AssertSuccessfulBehavior := func() {
		It("should not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should select area", func() {
			Expect(g.selectedAreaID).ToNot(BeZero())
		})
	}

	JustBeforeEach(func() {
		g, err = g.SelectThief(c)
	})

	// 	Describe("Select Thief", func() {
	// 		BeforeEach(func() {
	// 			c.Request = httptest.NewRequest(http.MethodPost, gamePath+selectThiefPath+"/1", nil)
	// 			c.Params = gin.Params{gin.Param{"hid", "1"}}
	// 			err = putHeader(ctx, 1)
	// 			Expect(err).To(BeNil())
	// 			_, err = header.GetHeader(ctx, 1)
	// 			Expect(err).To(BeNil())
	// 		})
	Context("when no current user", func() {
		It("should indicate there is no current user", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("only the current player can perform the selected action"))
		})

		Context("when there is a valid request", func() {

			BeforeEach(func() {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+selectThiefPath,
					strings.NewReader(`{ "row": 1 , "column": 1 }`),
				)
				c.Request.Header.Set("Content-Type", "application/json")
			})

			AssertFailedBehavior()

		})
	})
	// Context("without current user", func() {
	// 	It("should warn of missing current user", func() {
	// 		// selectThief()(c)
	// 		result := resp.Body.String()
	// 		Expect(result).To(ContainSubstring("current user not found"))
	// 	})
	// })
	Context("when the current user is the current player", func() {

		BeforeEach(func() {
			if g.CPUserIndices[0] == 1 {
				cu = u1
			} else {
				cu = u2
			}
			user.WithCurrent(c, cu)

			cp = g.currentPlayerFor(cu)
			a = g.grid.area(1, 1)
			g.grid, a = g.grid.placeThiefIn(cp, a)
		})

		Context("when there is a valid request", func() {

			BeforeEach(func() {
				c.Request = httptest.NewRequest(
					http.MethodPost,
					"/"+selectThiefPath,
					strings.NewReader(`{ "row": 1 , "column": 1 }`),
				)
				c.Request.Header.Set("Content-Type", "application/json")
			})

			AssertSuccessfulBehavior()

		})
		// Context("with current user", func() {
		// 	BeforeEach(func() {
		// 		Expect(u).ToNot(BeNil())
		// 		user.WithCurrent(c, u)
		// 	})
		// 	It("should not warn of missing current user", func() {
		// 		selectThief()(c)
		// 		result := resp.Body.String()
		// 		Expect(result).ToNot(ContainSubstring("Current user not found."))
		// 	})
		// 	Context("with correct json params", func() {
		// 		It("should create game", func() {
		// 			selectThief()(c)
		// 		})
		// 	})
		// 	Context("With incorrect json params", func() {
		// 		It("should not create game", func() {
		// 			selectThief()(c)
		// 		})
		// 	})
		// })
	})
})
