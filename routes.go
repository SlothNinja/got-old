package main

import (
	"github.com/SlothNinja/sn"
	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
)

const (
	authPath = "/auth"
	// mailPrefix  = "/mail"
	hidParam    = "hid"
	countParam  = "count"
	offsetParam = "offset"
	statusParam = "status"
	uidParam    = "uid"
	endQueue    = "end"

	devRoot         = "localhost:8081"
	prodRoot        = "got.slothninja.com"
	startQueue      = "start"
	gamePath        = "game"
	gamesPath       = gamePath + "s"
	showPath        = gamePath + "/show"
	undoPath        = gamePath + "/undo"
	resetPath       = gamePath + "/reset"
	finishPath      = gamePath + "/finish"
	endPath         = gamePath + "/end"
	redoPath        = gamePath + "/redo"
	placeThiefPath  = gamePath + "/place-thief"
	playCardPath    = gamePath + "/play-card"
	selectThiefPath = gamePath + "/select-thief"
	moveThiefPath   = gamePath + "/move-thief"
	passPath        = gamePath + "/pass"
	addMessagePath  = gamePath + "/add-message"
	messagesPath    = gamePath + "/messages"
	glogPath        = gamePath + "/glog"
	newPath         = gamePath + "/new"
	createPath      = newPath
	startPath       = gamePath + "/start"
	dropPath        = gamePath + "/drop"
	acceptPath      = gamePath + "/accept"
	adminPath       = gamePath + "/admin"
)

func addRoutes(prefix string, engine *gin.Engine, s server) {
	// Mail route
	// mail.AddRoutes(mailPrefix, engine)

	// Group
	g1 := engine.Group(prefix)

	// Current User
	g1.GET("/current", current(prefix))

	// Loging
	g1.GET("/login", user.Login(authPath))

	// authHandler
	g1.GET("/auth", user.Auth(authPath))

	// New
	g1.GET(newPath,
		s.newInvitation(),
	)

	// Create
	g1.PUT(createPath,
		s.createInvitation(),
	)

	// Show
	g1.GET(showPath+"/:"+hidParam,
		s.show(),
	)

	// Drop
	g1.PUT(dropPath+"/:"+hidParam,
		s.dropInvitation(hidParam),
	)

	// Accept
	g1.PUT(acceptPath+"/:"+hidParam,
		s.acceptInvitation(hidParam),
	)

	// Index
	g1.GET(gamesPath+"/:"+statusParam,
		s.jsonIndexAction(),
	)

	// JSON Data for Index
	g1.POST(gamesPath+"/:"+statusParam+"/json",
		s.jsonIndexAction(),
	)

	// JSON Data for Index
	g1.POST(gamesPath+"/:"+statusParam+"/user/:"+uidParam+"/json",
		s.jsonIndexAction(),
	)

	// Undo
	g1.PUT(undoPath+"/:"+hidParam,
		s.undo(hidParam),
	)

	// Redo
	g1.PUT(redoPath+"/:"+hidParam,
		s.redo(hidParam),
	)

	// Reset
	g1.PUT(resetPath+"/:"+hidParam,
		s.reset(hidParam),
	)

	// Finish
	g1.PUT(finishPath+"/:"+hidParam,
		s.finish(hidParam, endPath),
	)

	// Place Thief
	g1.PUT(placeThiefPath+"/:"+hidParam,
		s.placeThief(),
	)

	// Play Card
	g1.PUT(playCardPath+"/:"+hidParam,
		s.playCard(),
	)

	// Select Thief
	g1.PUT(selectThiefPath+"/:"+hidParam,
		s.selectThief(),
	)

	// Move Thief
	g1.PUT(moveThiefPath+"/:"+hidParam,
		s.moveThief(),
	)

	// Pass
	g1.PUT(passPath+"/:"+hidParam,
		s.pass(),
	)

	// Add Message
	g1.PUT(addMessagePath+"/:"+hidParam,
		sn.AddMLogMessage(),
	)

	// Get Messages
	g1.GET(messagesPath+"/:"+hidParam+"/:"+offsetParam,
		sn.GetMLog(hidParam, offsetParam),
	)

	// Get game Log
	g1.GET(glogPath+"/:"+hidParam+"/:"+countParam+"/:"+offsetParam,
		s.getGLog(hidParam, countParam, offsetParam),
	)
}
