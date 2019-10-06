package main

import "github.com/SlothNinja/sn"

const (
	phaseNone sn.Phase = iota
	phaseSetup
	phaseStartGame
	phasePlaceThieves
	phasePlayCard
	phaseSelectThief
	phaseMoveThief
	phaseClaimItem
	phaseFinalClaim
	phaseAnnounceWinners
	phaseGameOver
	phaseEndGame
)

func phaseName(p sn.Phase) string {
	switch p {
	case phaseSetup:
		return "phaseSetup"
	case phaseStartGame:
		return "Start game"
	case phasePlaceThieves:
		return "Place Thieves"
	case phasePlayCard:
		return "Play Card"
	case phaseSelectThief:
		return "Select Thief"
	case phaseMoveThief:
		return "Move Thief"
	case phaseClaimItem:
		return "Claim Magical Item"
	case phaseFinalClaim:
		return "Final Claim"
	case phaseAnnounceWinners:
		return "Announce Winners"
	case phaseGameOver:
		return "game Over"
	case phaseEndGame:
		return "End Of game"
	default:
		return "None"
	}
}
