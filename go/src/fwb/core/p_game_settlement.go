package core

import (
	"er"
)

func pgstInit(me *gameImp) *er.Err {
	me.lg.Dbg("Enter Game Settlement phase")

	return me.gotoPhase(_P_GAME_FINISH)
}
