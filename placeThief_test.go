package main

import (
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
			"/"+placeThiefPath+"/1",
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

var _ = Describe("g.placeThief", func() {
	var (
		c          *gin.Context
		cp         player
		a          area
		g          game
		cu, u1, u2 user.User
		err        error
	)

	BeforeEach(func() {
		c, _ = gin.CreateTestContext(httptest.NewRecorder())

		u1, u2 = createUsers()

		g = createGame(c, u1, u2)
	})

	JustBeforeEach(func() {
		g, err = g.placeThief(c)
	})

	Context("when no current user", func() {
		It("should indicate there is no current user", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("only the current player can perform the selected action"))
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
				a = g.grid.area(1, 1)
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

			cp = g.currentPlayerFor(cu)
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
				a = g.grid.area(1, 1)
				Expect(a.thief.pid).Should(Equal(cp.id))
			})

			Context("when thief already in selected area", func() {
				BeforeEach(func() {
					a = g.grid.area(1, 1)

					g.grid, a = g.grid.placeThiefIn(cp, a)
				})

				It("should indicate area already has thief", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when no card in selected area", func() {
				BeforeEach(func() {
					a = g.grid.area(1, 1)

					g.grid, a = g.grid.removeCardFrom(a)
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
				a = g.grid.area(1, 1)
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
				cp = g.currentPlayerFor(u2)
			} else {
				user.WithCurrent(c, u1)
				cp = g.currentPlayerFor(u1)
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
				Expect(err.Error()).To(ContainSubstring("only the current player can perform the selected action"))
			})

			It("should not place thief", func() {
				a = g.grid.area(1, 1)
				Expect(a.hasThief()).To(BeFalse())
			})
		})
	})
})
