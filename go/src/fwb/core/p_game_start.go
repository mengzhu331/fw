package core

import (
	"er"
	"fwb"
	"sgs"
)

type playerAckMap map[int]bool

type pgsData struct {
	pam playerAckMap
}

func pgsInit(me *gameImp) *er.Err {
	me.gd.Init(me.app.GetPlayers(), me.conf.MaxPawn, me.conf.MinRounds)

	me.setDCE(fwb.CMD_GAME_START_ACK, pgsOnGameStartAck)
	me.setTimer(10000, pgsOnTimeOut)
	me.pd = &pgsData{
		pam: make(playerAckMap),
	}
	return me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_GAME_START,
		Source:  fwb.CMD_SOURCE_APP,
		Payload: makePlayersInfo(me),
	})
}

func pgsOnGameStartAck(me *gameImp, command sgs.Command) *er.Err {
	pd := me.pd.(pgsData)
	pd.pam[command.Source] = true
	npa := 0

	for range pd.pam {
		npa++
	}

	if npa == len(me.app.GetPlayers()) {
		me.gotoPhase(_P_ROUNDS_START)
		return nil
	}

	return nil
}

func pgsOnTimeOut(me *gameImp, command sgs.Command) (bool, *er.Err) {
	noResPlayers := make([]int, 0)
	players := me.app.GetPlayers()
	pd := me.pd.(pgsData)

	for cid := range players {
		if _, found := pd.pam[cid]; !found {
			noResPlayers = append(noResPlayers, cid)
		}
	}

	return false, me.GameOver(fwb.GAME_OVER_PLAYER_TIMEOUT, noResPlayers)
}

func makePlayersInfo(me *gameImp) map[string]int {
	playersInfo := make(map[string]int)
	for _, p := range me.app.GetPlayers() {
		playersInfo[p.Name()] = p.ID()
	}
	return playersInfo
}
