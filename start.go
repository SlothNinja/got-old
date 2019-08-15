package main

import (
	"bitbucket.org/SlothNinja/status"
)

const startgameID = "start-game"

// Start begins a Guild of Thieves game.
func (g game) start() game {
	g.Status = status.Running
	g = g.setupPhase()
	g = g.beginningOfPhaseReset()
	g.Phase = phaseStartGame
	g.Log = nil
	g.Log = append(g.Log, logEntry{
		"template": startgameID,
		"phase":    g.Phase,
		"turn":     g.Turn,
		"pids":     pids(g.players),
	})

	g.Turn++
	g.Phase = phasePlaceThieves
	g.Stack = g.Stack.Update().Commit()
	return g
}
