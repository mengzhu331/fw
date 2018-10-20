package core

import (
	"er"
	"fwb"
	"sgs"
	"sort"
)

func pgfInit(me *gameImp) *er.Err {
	me.lg.Dbg("Enter Game Finish phase")

	playerRematch := make(map[int]bool)
	me.pd = &playerRematch

	gameRank := makeGameRank(me)

	me.setDCE(fwb.CMD_REMATCH, pgfOnRematch)
	me.setTimer(1000, pgfOnTimeOut)

	printRank(me, gameRank)

	err := me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_GAME_FINISH,
		Who:     fwb.CMD_WHO_APP,
		Payload: gameRank,
	})

	return err
}

func pgfOnRematch(me *gameImp, command sgs.Command) *er.Err {
	if me.gd.GetPDIndex(command.Who) < 0 {
		return er.Throw(fwb.E_CMD_INVALID_CLIENT, er.EInfo{
			"details":  "rematch client ID illegal",
			"clientID": command.Who,
		}).To(me.lg)
	}

	playerRematch := me.pd.(*map[int]bool)
	(*playerRematch)[command.Who] = true
	if len(*playerRematch) == len(me.gd.PData) {
		return me.gotoPhase(_P_GAME_START)
	}
	return nil
}

func pgfOnTimeOut(me *gameImp, command sgs.Command) *er.Err {
	return me.GameOver(fwb.GAME_OVER_NORMAL, nil)
}

type pdSorter []fwb.PlayerData

// Len is part of sort.Interface.
func (me *pdSorter) Len() int {
	return len(*me)
}

// Swap is part of sort.Interface.
func (me *pdSorter) Swap(i, j int) {
	(*me)[i], (*me)[j] = (*me)[j], (*me)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (me *pdSorter) Less(i, j int) bool {
	return (*me)[i][fwb.PD_PT_HEART] <= (*me)[j][fwb.PD_PT_HEART]
}

func makeGameRank(me *gameImp) map[int]int {
	pdsrt := pdSorter(me.gd.PData)

	sort.Sort(sort.Reverse(&pdsrt))

	gameRank := make(map[int]int)

	var rank = 1
	for i, pd := range pdsrt {
		if i > 0 && pd[fwb.PD_PT_HEART] < pdsrt[i-1][fwb.PD_PT_HEART] {
			rank++
		}
		pid := pd[fwb.PD_CLIENT_ID]
		gameRank[pid] = rank
	}

	return gameRank
}

func printRank(me *gameImp, rank map[int]int) {
	me.alg.Inf("Game Finish, player rank:")
	for k, v := range rank {
		me.alg.Inf("  Player %v, Rank %v", me.app.GetPlayer(k).Name(), v)
	}
}
