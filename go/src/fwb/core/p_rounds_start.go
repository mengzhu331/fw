package core

import (
	"er"
	"fwb"
	"sgs"
)

type prsData struct {
	pam playerAckMap
}

func prsInit(me *gameImp) *er.Err {
	me.pd = &prsData{
		pam: make(playerAckMap),
	}

	me.setDCE(fwb.CMD_ROUND_START_ACK, prsOnRoundStartAck)
	me.setTimer(10000, prsOnTimeOut)
	setTurnOrder(me)

	basicCards, shuffledCards := me.cm.MakeCardSet()
	cards := append(basicCards, shuffledCards...)
	return me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_ROUND_START,
		Source:  fwb.CMD_SOURCE_APP,
		Payload: cards,
	})
}

func setTurnOrder(me *gameImp) {
	me.turnOrder = make([]int, len(me.app.GetPlayers()))

	for id := range me.app.GetPlayers() {
		me.turnOrder = append(me.turnOrder, id)
	}

	minHeart := 9999
	minGold := 9999
	minIdx := 0

	for i, id := range me.turnOrder {
		heart := me.gd.PData[id][fwb.PD_PT_HEART]
		gold := me.gd.PData[id][fwb.PD_PT_GOLD]
		if heart < minHeart || (heart == minHeart && gold < minGold) {
			minGold = gold
			minHeart = heart
			minIdx = i
		}
	}

	me.turnOrder = append(me.turnOrder[minIdx:], me.turnOrder[:minIdx]...)
}

func prsOnRoundStartAck(me *gameImp, command sgs.Command) (bool, *er.Err) {
	pd := me.pd.(*prsData)
	pd.pam[command.Source] = true
	npa := 0

	for range pd.pam {
		npa++
	}

	if npa == len(me.app.GetPlayers()) {
		me.gotoPhase(_P_ROUNDS_TURNS)
		return false, nil
	}

	return true, nil
}

func prsOnTimeOut(me *gameImp, command sgs.Command) (bool, *er.Err) {
	//we continue if timeout
	me.gotoPhase(_P_ROUNDS_TURNS)
	return false, nil
}
