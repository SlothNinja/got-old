package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/user"
)

func (g game) validatePlayerAction(c *gin.Context) (player, error) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cu, found := user.Current(c)
	if !found {
		return player{}, errors.WithMessage(errValidation, "unable to find current user")
	}
	cp, found := g.currentPlayerFor(cu)
	switch {
	case !found:
		return player{}, errors.WithMessage(errValidation, "current player not found")
	case !g.CPorAdmin(cp.id, cu):
		return player{}, errors.WithMessage(errValidation, "only the current player can perform the selected action")
	case cp.performedAction:
		return player{}, errors.WithMessage(errValidation, "you have already performed an action")
	}
	return cp, nil
}

func (g game) validateAdminAction(c *gin.Context) error {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	cu, found := user.Current(c)
	if !found {
		return errors.WithMessage(errValidation, "unable to find current user")
	}
	if !cu.Admin {
		return errors.WithMessage(errValidation, "only an admin can perform the selected action")
	}
	return nil
}
