package core

import (
	"er"
	"fwb"
	"sgs"

	"github.com/google/uuid"
)

type playerAckMap map[int]bool

type pgsData struct {
	pam playerAckMap
}

func pgsInit(me *gameImp) *er.Err {
	me.lg.Dbg("Enter Game Start phase")

	me.gameuuid, _ = uuid.NewUUID()
	me.alg = _gameLog.Child("Game_" + me.gameuuid.String())

	me.gd.Init(me.app.GetPlayers(), me.conf.MaxPawn, me.conf.MinRounds, me.conf.StartGold)

	me.setDCE(fwb.CMD_GAME_START_ACK, pgsOnGameStartAck)
	me.setTimer(10000, pgsOnTimeOut)
	me.pd = &pgsData{
		pam: make(playerAckMap),
	}
	return me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_GAME_START,
		Who:     fwb.CMD_WHO_APP,
		Payload: makePlayersInfo(me),
	})
}

func pgsOnGameStartAck(me *gameImp, command sgs.Command) *er.Err {
	pd := me.pd.(*pgsData)
	pd.pam[command.Who] = true
	applyBasicSkill(me, command)

	npa := 0

	err := me.app.SendToPlayer(command.Who, sgs.Command{
		ID:      fwb.CMD_SYNC_GAME_STATE,
		Who:     fwb.CMD_WHO_APP,
		Payload: me.gd,
	})

	if err.Importance() >= er.IMPT_DEGRADE {
		return err
	}

	for range pd.pam {
		npa++
	}

	if npa == len(me.app.GetPlayers()) {
		me.gotoPhase(_P_ROUNDS_START)
	}

	return nil
}

func pgsOnTimeOut(me *gameImp, command sgs.Command) *er.Err {
	noResPlayers := make([]int, 0)
	players := me.app.GetPlayers()
	pd := me.pd.(*pgsData)

	for _, p := range players {
		if _, found := pd.pam[p.ID()]; !found {
			noResPlayers = append(noResPlayers, p.ID())
		}
	}

	return me.GameOver(fwb.GAME_OVER_PLAYER_TIMEOUT, noResPlayers)
}

func makePlayersInfo(me *gameImp) map[string]int {
	playersInfo := make(map[string]int)
	for _, p := range me.app.GetPlayers() {
		playersInfo[p.Name()] = p.ID()
	}
	return playersInfo
}

func applyBasicSkill(me *gameImp, command sgs.Command) *er.Err {
	me.lg.Inf("Apply game start skill %v %v", command.Payload, command.Who)
	p := me.gd.GetPDIndex(command.Who)
	if p < 0 {
		return er.Throw(fwb.E_CMD_INVALID_CLIENT, er.EInfo{
			"details": "game start ack from invalid client",
			"ID":      command.Who,
		}).To(me.lg)
	}

	switch command.Payload {
	case "inte":
		me.gd.PData[p][fwb.PD_SK_INTELLIGENCE]++
	case "know":
		me.gd.PData[p][fwb.PD_SK_KNOWLEDGE]++
	case "stre":
		me.gd.PData[p][fwb.PD_SK_STRENGTH]++
	default:
		return er.Throw(fwb.E_CMD_INVALID_PAYLOAD, er.EInfo{
			"details": "unknown payload in game start ack",
			"payload": command.Payload,
		}).To(me.lg)
	}

	return nil
}
