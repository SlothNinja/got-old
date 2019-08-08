package main

import (
	"github.com/gin-gonic/gin"

	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/mail"
)

const (
	path  = "/got/mail"
	queue = "mail"
)

func sendTurnNotification(c *gin.Context, h Header, pid int) error {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	// var m *mail.Message

	// //	if h.NotificationsByPID(pid) {
	// m = noticeOfTurn(c, h, pid)
	// err = send(c, m)
	// //	}
	return nil
}

func sendEndGameNotifications(c *gin.Context, h Header, rs []result) error {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	// var ms []*mail.Message

	// for i := range h.Users {
	// 	pid := i + 1
	// 	ms = append(ms, noticeOfEndGame(c, h, rs, pid))
	// }

	// if len(ms) > 1 {
	// 	err = send(c, ms...)
	// }
	return nil
}

// func send(c context.Context, ms ...*mail.Message) error {
func send(c *gin.Context) error {
	//return mail.send(c, path, queue, ms...)
	return mail.Send(c, path, queue)
}
