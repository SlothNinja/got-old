package main

import (
	"bitbucket.org/SlothNinja/status"
	"bitbucket.org/SlothNinja/user"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Got Suite")
}

func createUsers() (user.User2, user.User2) {
	u1 := user.New2("1")
	u1.Name = "SlothNinja1"

	u2 := user.New2("2")
	u2.Name = "steve2"

	return u1, u2
}

func createGame(c *gin.Context, u1, u2 user.User2) game {
	h := newHeaderEntity(newGame(1))
	h.Creator = u1
	h.Title = "SlothNinja's Game"
	h.NumPlayers = 2

	h, start, err := h.Accept(c, u1)
	Expect(err).NotTo(HaveOccurred())
	Expect(start).To(BeFalse())

	h, start, err = h.Accept(c, u2)
	Expect(err).NotTo(HaveOccurred())
	Expect(start).To(BeTrue())

	h.Status = status.Starting

	g := newGame(h.ID())
	g.Header = h.Header
	g = g.start()
	return g
}

func createUsers3() (user.User2, user.User2, user.User2) {
	u1 := user.New2("1")
	u1.Name = "SlothNinja1"

	u2 := user.New2("2")
	u2.Name = "steve2"

	u3 := user.New2("3")
	u3.Name = "george3"

	return u1, u2, u3
}

func createGame3(c *gin.Context, u1, u2, u3 user.User2) game {
	h := newHeaderEntity(newGame(1))
	h.Creator = u1
	h.Title = "SlothNinja's Game"
	h.NumPlayers = 3

	h, start, err := h.Accept(c, u1)
	Expect(err).NotTo(HaveOccurred())
	Expect(start).To(BeFalse())

	h, start, err = h.Accept(c, u2)
	Expect(err).NotTo(HaveOccurred())
	Expect(start).To(BeFalse())

	h, start, err = h.Accept(c, u3)
	Expect(err).NotTo(HaveOccurred())
	Expect(start).To(BeTrue())

	h.Status = status.Starting

	g := newGame(h.ID())
	g.Header = h.Header
	g = g.start()
	return g
}
