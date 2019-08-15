package main

import (
	"fmt"
	"sort"

	"bitbucket.org/SlothNinja/contest"
	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/rating"
	"bitbucket.org/SlothNinja/status"
	"github.com/gin-gonic/gin"
)

// func (g game) Finishgame(c *gin.Context, u user.User) error {
//
// 	rs, err := g.endgame(c)
// 	if err != nil {
// 		return err
// 	}
//
// 	// if err = g.Commit(ctx, u); err != nil {
// 	// 	return
// 	// }
//
// 	// Log but, otherwise ignore error.  game will finish even if unable to mail end game results.
// 	err = sendEndGameNotifications(c, g.Header, rs)
// 	return err
// }

func (g game) endgame(c *gin.Context) ([]result, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	g.finalClaim(c)

	ps, err := g.determinePlaces(c)
	if err != nil {
		return nil, err
	}

	g.setWinners(ps[0])

	// Log winners
	// g.GLog.AddEntryData(glog.EntryData{
	// 	"template": "announce-winners",
	// 	"turn":     g.Turn,
	// 	"phase":    g.Phase,
	// 	"winners":  g.WinnerIndices,
	// })

	// // Log end game
	// g.GLog.AddEntryData(glog.EntryData{
	// 	"template": "end-game",
	// 	"turn":     g.Turn,
	// 	"phase":    g.Phase,
	// })

	cs := contest.GenContests(c, ps)

	g.Status = status.Completed

	// Need to call endgameResults before saving the new contests.
	// endgameResults relies on pulling the old contests from the datastore.
	// Saving the contests results in double counting.
	return g.endgameResults(c, ps, cs)
}

func (g game) toIDS(places [][]player) [][]string {
	sids := make([][]string, len(places))
	for i, ps := range places {
		for _, p := range ps {
			sids[i] = append(sids[i], g.UserIDByPID(p.id))
		}
	}
	return sids
}

//type endgameEntry struct {
//	Entry
//}
//
//func (g game) newEndgameEntry() {
//	e := &endgameEntry{
//		Entry: g.newEntry(),
//	}
//	g.Log = append(g.Log, e)
//}
//
//func (e endgameEntry) HTML(g game) (s template.HTML) {
//	rows := sn.HTML("")
//	for _, p := range g.Players() {
//		rows += sn.HTML("<tr>")
//		rows += sn.HTML("<td>%s</td> <td>%d</td> <td>%d</td> <td>%d</td> <td>%d</td>",
//			g.NameFor(p), p.score, lampCount(p.Hand...), camelCount(p.Hand...), len(p.Hand))
//		rows += sn.HTML("</tr>")
//	}
//	s += sn.HTML("<table class='strippedDataTable'><thead><tr><th>Player</th><th>score</th>")
//	s += sn.HTML("<th>Lamps</th><th>Camels</th><th>Cards</th></tr></thead><tbody>")
//	s += rows
//	s += sn.HTML("</tbody></table>")
//	return
//}

func (g game) setWinners(rmap contest.ResultsMap) {
	g.Phase = phaseAnnounceWinners
	g.Status = status.Completed

	g.CPUserIndices = nil
	g.WinnerIndices = nil
	for k := range rmap {
		pid, _ := g.PlayerIDFor(k.Name)
		g.WinnerIndices = append(g.WinnerIndices, pid)
	}

}

func (g game) endgameResults(c *gin.Context, ps contest.Places, cs contest.Contests) ([]result, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	rs := make([]result, g.NumPlayers)
	i := 0

	for place, rmap := range ps {
		for k := range rmap {
			pid, found := g.PlayerIDFor(k.Name)
			if !found {
				return nil, fmt.Errorf("player for uid: %v not found", k.Name)
			}
			p, found := playerByID(pid, g.players)
			if !found {
				return nil, fmt.Errorf("player with pid: %v not found", pid)
			}
			cr, nr, err := rating.IncreaseFor(c, k, g.Type, cs)
			if err != nil {
				return nil, err
			}
			clo, nlo := cr.Rank().GLO(), nr.Rank().GLO()
			inc := nlo - clo

			rs[i] = result{
				Place: place + 1,
				GLO:   nlo,
				Score: p.score,
				Name:  g.NameByPID(pid),
				Inc:   fmt.Sprintf("%+d", inc),
			}
		}
		i++
	}

	return rs, nil
}

//type announceWinnersEntry struct {
//	Entry
//}
//
//func (g game) newAnnounceWinnersEntry() announceWinnersEntry {
//	e := &announceWinnersEntry{
//		Entry: g.newEntry(),
//	}
//	g.Log = append(g.Log, e)
//	return e
//}
//
//func (e announceWinnersEntry) HTML(g game) template.HTML {
//	names := make([]string, len(g.winners()))
//	for i, winner := range g.winners() {
//		names[i] = g.NameFor(winner)
//	}
//	return sn.HTML("Congratulations: %s.", restful.toSentence(names))
//}

func (g game) winners() []player {
	l := len(g.WinnerIndices)
	if l == 0 {
		return nil
	}
	ps := make([]player, l)
	for i, pid := range g.WinnerIndices {
		p, found := playerByID(pid, g.players)
		if found {
			ps[i] = p
		}
	}
	return ps
}

func (g game) determinePlaces(c *gin.Context) (contest.Places, error) {
	// sort players by score with greatest first
	ps := g.players
	sort.SliceStable(ps, func(i, j int) bool {
		b, _ := ps[i].defeated(ps[j])
		return b
	})
	g.players = ps

	uids := playerUKeys(g.players)
	rs, err := rating.GetMulti(c, uids, g.Type)
	if err != nil {
		return nil, err
	}

	places := make(contest.Places, 0)
	rmap := make(contest.ResultsMap, 0)
	for i, p1 := range g.players {
		results := make(contest.Results, 0)
		tie := false
		for j, p2 := range g.players {
			p2Rating := rs[j]
			result := &contest.Result{
				GameID: g.ID(),
				R:      p2Rating.R,
				RD:     p2Rating.RD,
			}
			switch b, tie := p1.defeated(p2); {
			case i == j:
			case tie == true:
				result.Outcome = 0.5
			case b == true:
				result.Outcome = 1
			default:
				result.Outcome = 0
			}
			results = append(results, result)
		}
		rmap[uids[i]] = results
		if !tie {
			places = append(places, rmap)
			rmap = make(contest.ResultsMap, 0)
		} else if i == len(g.players)-1 {
			places = append(places, rmap)
		}
	}
	return places, nil
}
