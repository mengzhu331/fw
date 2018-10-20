package core

import (
	"er"
	"fwb"
	"sgs"
	"strconv"
)

type prsData struct {
	pam playerAckMap
}

func prsInit(me *gameImp) *er.Err {
	me.lg.Dbg("Enter Round Start phase")

	me.gd.Round++
	for _, p := range me.gd.PData {
		p[fwb.PD_PAWNS] = p[fwb.PD_MAX_PAWNS]
	}

	me.pd = &prsData{
		pam: make(playerAckMap),
	}

	me.setDCE(fwb.CMD_ROUND_START_ACK, prsOnRoundStartAck)
	me.setTimer(2000, prsOnTimeOut)
	setTurnOrder(me)

	specialCards, basicCards, shuffledCards := me.cm.MakeCardSet()
	cards := append(append(specialCards, basicCards...), shuffledCards...)
	me.gd.Cards = cards

	printRoundInfo(me)

	return me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_ROUND_START,
		Who:     fwb.CMD_WHO_APP,
		Payload: cards,
	})
}

func setTurnOrder(me *gameImp) {
	me.turnOrder = make([]int, 0, len(me.app.GetPlayers()))

	for _, p := range me.app.GetPlayers() {
		me.turnOrder = append(me.turnOrder, p.ID())
	}

	minHeart := 9999
	minGold := 9999
	minIdx := 0

	for i, id := range me.turnOrder {
		px := me.gd.GetPDIndex(id)
		heart := me.gd.PData[px][fwb.PD_PT_HEART]
		gold := me.gd.PData[px][fwb.PD_PT_GOLD]
		if heart < minHeart || (heart == minHeart && gold < minGold) {
			minGold = gold
			minHeart = heart
			minIdx = i
		}
	}

	me.turnOrder = append(me.turnOrder[minIdx:], me.turnOrder[:minIdx]...)
}

func printRoundInfo(me *gameImp) {
	cardsList := "["
	for i, c := range me.gd.Cards {
		cardsList += strconv.Itoa(c.ID)
		if i == len(me.gd.Cards)-1 {
			break
		}
		cardsList += ", "
	}
	cardsList += "]"

	me.alg.Inf("Round %v, current cards %v", me.gd.Round, cardsList)
	me.alg.Inf("Player Data")
	for _, p := range me.gd.PData {
		me.alg.Inf("  Player %v: %v", me.app.GetPlayer(p[fwb.PD_CLIENT_ID]).Name(), p[fwb.PD_CLIENT_ID+1:])
	}
}

func prsOnRoundStartAck(me *gameImp, command sgs.Command) *er.Err {
	pd := me.pd.(*prsData)
	pd.pam[command.Who] = true
	npa := 0

	for range pd.pam {
		npa++
	}

	if npa == len(me.app.GetPlayers()) {
		me.gotoPhase(_P_ROUNDS_TURNS)
		return nil
	}

	return nil
}

func prsOnTimeOut(me *gameImp, command sgs.Command) *er.Err {
	//we continue if timeout
	me.gotoPhase(_P_ROUNDS_TURNS)
	return nil
}
