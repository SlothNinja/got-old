package main

import (
	"fmt"
	"strings"

	"bitbucket.org/SlothNinja/log"
	"github.com/gin-gonic/gin"

	"go.chromium.org/gae/service/mail"
)

func noticeOfTurn(c *gin.Context, g game, to int) (m mail.Message) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	subject := fmt.Sprintf("It's your turn in %s (%s #%d).", g.Type, g.Title, g.ID())
	url := fmt.Sprintf(`<a href="http://www.slothninja.com/got/game/show/%d">here</a>`, g.ID())
	body := fmt.Sprintf(`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
		<html>
			<head>
				<meta http-equiv="content-type" content="text/html; charset=ISO-8859-1">
			</head>
			<body bgcolor="#ffffff" text="#000000">
				<p>%s</p>
				<p>You can take your turn %s.</p>
			</body>
		</html>`, subject, url)

	m = mail.Message{
		Sender:   "webmaster@slothninja.com",
		To:       []string{g.EmailByPID(to)},
		Subject:  subject,
		HTMLBody: body,
	}

	return
}

func noticeOfEndGame(c *gin.Context, g game, rs []result, to int) mail.Message {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	var names []string
	for _, pid := range g.WinnerIndices {
		names = append(names, g.NameByPID(pid))
	}

	body := `
	  <!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
	  <html>
		  <head>
			  <meta http-equiv="content-type" content="text/html; charset=ISO-8859-1">
		  </head>
		  <body bgcolor="#ffffff" text="#000000">`
	for _, r := range rs {
		body += fmt.Sprintf(`
			  <div style="height:3em">
				  <div style="height:3em;float:left;padding-right:1em">%d.</div>
				  <div style="height:1em">%s scored %d points.</div>
				  <div style="height:1em">Glicko %s (-> %d)</div>
			  </div>`, r.Place, r.Name, r.Score, r.Inc, r.GLO)
	}
	body += fmt.Sprintf(`
			  <p></p>
			  <p>Congratulations: %s.</p>
		  </body>
	  </html>`, toSentence(names))

	sender := "webmaster@slothninja.com"
	subject := fmt.Sprintf("SlothNinja games: Guild of Thieves #%d Has Ended", g.ID())
	return mail.Message{
		To:       []string{g.EmailByPID(to)},
		Sender:   sender,
		Subject:  subject,
		HTMLBody: body,
	}
}

func toSentence(ss []string) (s string) {
	switch l := len(ss); l {
	case 0:
	case 1:
		s = ss[0]
	case 2:
		s = ss[0] + " and " + ss[1]
	default:
		s = strings.Join(ss[:l-1], ", ") + ", and " + ss[l-1]
	}
	return s
}
