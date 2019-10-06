package main

import (
	"fmt"

	"github.com/SlothNinja/log"
	"github.com/SlothNinja/user"

	"github.com/gin-gonic/gin"
)

func (g game) Pass(c *gin.Context) (game, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp, err := g.validatePass(c)
	if err != nil {
		return g, err
	}

	cp.passed = true
	cp.performedAction = true

	g.Phase = phaseClaimItem
	cu := user.Current(c)
	g.updateClickablesFor(cu)

	// Log Pass
	// g.GLog.SetEntryData(glog.EntryData{
	// 	"template": "pass",
	// 	"turn":     g.Turn,
	// 	"phase":    g.Phase,
	// 	"pid":      cp.ID,
	// })

	return g, nil
}

func (g game) validatePass(c *gin.Context) (player, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	// if cp is nil, then err will not be nil.
	cp, err := g.validatePlayerAction(c)
	switch {
	case err != nil:
		return noPlayer, err
	case g.Phase != phasePlayCard:
		return noPlayer, fmt.Errorf("wrong phase for selected action: %w", errValidation)
	}
	return cp, nil
}

//type passEntry struct {
//	*Entry
//}
//
//func (g *game) newPassEntryFor(p player) (e *passEntry) {
//	e = &passEntry{
//		Entry: g.newEntryFor(p),
//	}
//	p.Log = append(p.Log, e)
//	g.Log = append(g.Log, e)
//	return
//}
//
//func (e *passEntry) HTML(g *game) template.HTML {
//	return sn.HTML("%s passed.", g.NameByPID(e.PlayerID))
//}
