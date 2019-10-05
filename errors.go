package main

import (
	"errors"
	"fmt"
)

var (
	errValidation         = errors.New("validation error")
	errUnexpected         = errors.New("unexpected error")
	errUserNotFound       = fmt.Errorf("current user not found: %w", errValidation)
	errPlayerNotFound     = fmt.Errorf("player not found: %w", errValidation)
	errActionNotPerformed = fmt.Errorf("player has yet to perform an action: %w", errValidation)
	errNotCPorAdmin       = fmt.Errorf("not current player or admin: %w", errValidation)
	errWrongPhase         = fmt.Errorf("wrong phase: %w", errValidation)
	errMissingToken       = errors.New("missing token")
)
