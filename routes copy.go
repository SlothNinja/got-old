package main

import (
	"bitbucket.org/SlothNinja/chat"
	"github.com/gin-gonic/gin"
)

const (
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

// AddRoutes adds routes to engine.
func AddRoutes(engine *gin.Engine, s server) {
	g1 := engine.Group("")

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
		jsonIndexAction(),
	)

	//g1.GET(gamesPath+"/:"+statusParam+"/user/:"+uidParam,
	//	Index(),
	//)

	// JSON Data for Index
	g1.POST(gamesPath+"/:"+statusParam+"/json",
		jsonIndexAction(),
	)

	// JSON Data for Index
	g1.POST(gamesPath+"/:"+statusParam+"/user/:"+uidParam+"/json",
		jsonIndexAction(),
	)
	// // Start
	// g1.POST(startPath+"/:"+hidParam,
	// 	s.start(),
	// )

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

	// JSON Index Data
	//g1.GET(gamesPath+"/:"+statusParam,
	//	s.jsonIndexAction(),
	//)

	//// // Index
	//// g1.GET(gamesPath+"/:"+statusParam+"/user/:"+uidParam,
	//// 	s.index(),
	//// )

	//// JSON Data for Index
	//g1.POST(gamesPath+"/:"+statusParam+"/json",
	//	s.jsonIndexAction(),
	//)

	//// JSON Data for Index
	//g1.POST(gamesPath+"/:"+statusParam+"/user/:"+uidParam+"/json",
	//	s.jsonIndexAction(),
	//)

	// Add Message
	g1.PUT(addMessagePath+"/:"+hidParam,
		chat.AddMLogMessage(),
	)

	// Get Messages
	g1.GET(messagesPath+"/:"+hidParam+"/:"+offsetParam,
		chat.GetMLog(hidParam, offsetParam),
	)

	// Get game Log
	g1.GET(glogPath+"/:"+hidParam+"/:"+countParam+"/:"+offsetParam,
		s.getGLog(hidParam, countParam, offsetParam),
	)
}
