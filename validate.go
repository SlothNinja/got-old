package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/errors/fmt"

	"github.com/SlothNinja/log"
	"github.com/SlothNinja/user"
)

func (g game) validatePlayerAction(c *gin.Context) (player, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cu := user.Current(c)
	cp := g.currentPlayerFor(cu)
	switch {
	case !g.CPorAdmin(cp.id, cu):
		return noPlayer, fmt.Errorf("only the current player can perform the selected action: %w", errValidation)
	case cp.performedAction:
		return noPlayer, fmt.Errorf("you have already performed an action: %w", errValidation)
	}
	return cp, nil
}

func (g game) validateAdminAction(c *gin.Context) error {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cu := user.Current(c)
	if !cu.Admin {
		return fmt.Errorf("only an admin can perform the selected action: %w", errValidation)
	}
	return nil
}
