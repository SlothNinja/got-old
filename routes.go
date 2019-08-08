package main

import (
	"bitbucket.org/SlothNinja/mail"
	"bitbucket.org/SlothNinja/user"
	"github.com/gin-gonic/gin"
)

const (
	authPath   = "/auth"
	mailPrefix = "/mail"
)

func addRoutes(prefix string, engine *gin.Engine, s server) {
	// Guild of Thieves Game Routes
	AddRoutes(engine, s)

	// Guild of Thieves Header Routes
	// hcon.AddRoutes(engine)

	// Mail route
	mail.AddRoutes(mailPrefix, engine)

	// User Group
	g1 := engine.Group(prefix)

	// New User
	// g1.GET("/new", newAction(prefix))

	// Create User
	// g1.PUT("/new", create(prefix))

	// Current User
	g1.GET("/current", current(prefix))

	// User
	// g1.GET("/json/:id", json(prefix))

	// Update User
	// g1.PUT("edit/:uid", update(prefix))

	g1.GET("/login", user.Login(authPath))

	// g1.GET("/logout", Logout)

	// authHandler
	g1.GET("/auth", user.Auth(authPath))

	// devauthHandler
	// g1.POST("/auth", DevAuth)
}
