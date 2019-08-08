package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"bitbucket.org/SlothNinja/user"
)

func (g game) validatePlayerAction(c *gin.Context) (player, error) {
	cu, found := user.Current(c)
	if !found {
		return player{}, fmt.Errorf("unable to find current user")
	}
	cp, found := g.currentPlayerFor(cu)
	switch {
	case !found:
		return player{}, errors.Wrap(errValidation, "current player not found")
	case !g.CPorAdmin(cp.ID, cu):
		return player{}, errors.Wrap(errValidation, "only the current player can perform the selected action")
	case cp.PerformedAction:
		return player{}, errors.Wrap(errValidation, "you have already performed an action")
	}
	return cp, nil
}

func (g game) validateAdminAction(c *gin.Context) error {
	cu, found := user.Current(c)
	if !found {
		return errors.WithMessage(errValidation, "unable to find current user")
	}
	if !cu.Admin {
		return errors.WithMessage(errValidation, "only an admin can perform the selected action")
	}
	return nil
}
