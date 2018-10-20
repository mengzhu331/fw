package core

import (
	"er"
	"math/rand"
)

func prfInit(me *gameImp) *er.Err {
	me.lg.Dbg("Enter Round Finish phase")

	if me.gd.Round >= me.gd.MinRounds {
		if rand.Intn(6)+1-me.gd.MinRounds > 3-me.gd.Round {
			return me.gotoPhase(_P_GAME_SETTLEMENT)
		}
	}

	return me.gotoPhase(_P_ROUNDS_START)
}
