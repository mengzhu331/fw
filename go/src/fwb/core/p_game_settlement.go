package core

import (
	"er"
)

func pgstInit(me *gameImp) *er.Err {
	return me.gotoPhase(_P_GAME_FINISH)
}
