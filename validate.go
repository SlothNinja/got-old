package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"
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
		return errors.WithMessage(errValidation, "only an admin can perform the selected action")
	}
	return nil
}
