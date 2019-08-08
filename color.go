package main

import (
	"bitbucket.org/SlothNinja/color"
)

var defaultColors = []color.Color{color.Yellow, color.Purple, color.Green, color.Black}

// func (g game) colorByPIDFor(u user.User2) func(int) color.Color {
//
// 	colors := defaultColors
//
// 	pid, found := g.PlayerIDFor(u.ID())
// 	if found {
// 		p, found := playerByID(pid, g.Players)
// 		if found {
// 			colors = p.Colors
// 		}
// 	}
//
// 	return func(pid int) (c color.Color) {
// 		i := pid - 1
// 		c = color.None
// 		if i >= 0 && i < len(colors) {
// 			c = colors[i]
// 		}
// 		return
// 	}
// }

// type colorizer interface {
// 	colorize(g game, u user.User2)
// }

// func (g game) UpdateColorsFor(u user.User2, ms []move.Move) {
// 	for i := range ms {
// 		if v, ok := ms[i].Data.(colorizer); ok {
// 			v.colorize(g, u)
// 		}
// 	}
// }
