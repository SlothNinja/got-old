package main

import "bitbucket.org/SlothNinja/log"

func (g game) setupPhase() game {
	g.Phase = phaseSetup
	g.OrderIndices = make([]int, g.NumPlayers)

	g = g.addNewPlayers()
	playerShuffle(g.players)
	g.grid = newGrid(g.NumPlayers)

	for i, p := range g.players {
		g.OrderIndices[i] = p.id
	}

	cp := g.players[len(g.players)-1]
	g.CPUserIndices = []int{cp.id}
	log.Debugf("g.CPUserIndices: %#v", g.CPUserIndices)
	return g
}
