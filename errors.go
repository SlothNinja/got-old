package main

import (
	"github.com/pkg/errors"
)

var (
	errValidation   = errors.New("validation error")
	errUnexpected   = errors.New("unexpected error")
	errUserNotFound = errors.New("current user not found")
)

// func jsonError(c *gin.Context, err error) {
// 	jsonErrorf(c, err.Error())
// }
//
// func jsonErrorf(c *gin.Context, format string, args ...interface{}) {
// 	log.Debugf(format, args...)
// 	msg := format
// 	if len(args) > 0 {
// 		msg = fmt.Sprintf(format, args...)
// 	}
// 	c.JSON(http.StatusOK, struct {
// 		Message string `json:"message"`
// 		Error   bool   `json:"error"`
// 	}{msg, true})
// }
