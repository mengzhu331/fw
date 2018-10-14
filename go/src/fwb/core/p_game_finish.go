package core

import (
	"er"
	"fwb"
	"sgs"
)

func pgfInit(me *gameImp) *er.Err {
	playerRematch := make(map[int]bool)
	me.pd = &playerRematch

	gameRank := makeGameRank(me)

	me.setDCE(fwb.CMD_REMATCH, pgfOnRematch)
	me.setTimer(10000, pgfOnTimeOut)

	err := me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_GAME_FINISH,
		Source:  fwb.CMD_SOURCE_APP,
		Payload: gameRank,
	})

	return err
}

func pgfOnRematch(me *gameImp, command sgs.Command) *er.Err {
	if me.gd.GetPDIndex(command.Source) < 0 {
		return er.Throw(fwb.E_CMD_INVALID_CLIENT, er.EInfo{
			"details":  "rematch client ID illegal",
			"clientID": command.Source,
		}).To(me.lg)
	}

	playerRematch := me.pd.(*map[int]bool)
	(*playerRematch)[command.Source] = true
	if len(*playerRematch) == len(me.gd.PData) {
		return me.gotoPhase(_P_GAME_START)
	}
	return nil
}

func pgfOnTimeOut(me *gameImp, command sgs.Command) (bool, *er.Err) {
	gameRank := makeGameRank(me)
	return false, me.GameOver(fwb.GAME_OVER_NORMAL, gameRank)
}

func makeGameRank(me *gameImp) map[int]int {
	pn := len(me.gd.PData)
	gameRank := make(map[int]int)
	sortedPlayers := make([]int, pn)

	for _, pd := range me.gd.PData {
		var i int
		for i = range sortedPlayers {
			if pd[fwb.PD_PT_HEART] >= me.gd.PData[sortedPlayers[i]][fwb.PD_PT_HEART] {
				break
			}
		}
		sortedPlayers = append(append(sortedPlayers[:i], pd[fwb.PD_CLIENT_ID]), sortedPlayers[i:]...)
	}

	var rank = 1
	for i, pid := range sortedPlayers {
		if i > 0 && me.gd.PData[pid][fwb.PD_PT_HEART] < me.gd.PData[sortedPlayers[i-1]][fwb.PD_PT_HEART] {
			rank++
		}
		gameRank[pid] = rank
	}

	return gameRank
}
