package core

import (
	"er"
	"fwb"
	"sgs"
)

type prtData struct {
	hotIndex int
	timer    int
}

func prtInit(me *gameImp) *er.Err {
	pd := &prtData{
		hotIndex: len(me.app.GetPlayers()) - 1,
		timer:    -1,
	}

	me.pd = pd

	return nextTurn(me)
}

func findNextPlayer(me *gameImp) int {
	pd := me.pd.(*prtData)
	pn := len(me.app.GetPlayers())
	pd.hotIndex = (pd.hotIndex + 1) % pn

	pid := me.turnOrder[pd.hotIndex]
	pindex := me.gd.GetPDIndex(pid)

	trivialHotIndex := pd.hotIndex

	for me.gd.PData[pindex][fwb.PD_PAWNS] <= 0 {
		pd.hotIndex = (pd.hotIndex + 1) % pn
		if pd.hotIndex == trivialHotIndex {
			return -1
		}

		pid = me.turnOrder[pd.hotIndex]
		pindex = me.gd.GetPDIndex(pid)
	}
	return pd.hotIndex
}

func nextTurn(me *gameImp) *er.Err {
	pd := me.pd.(*prtData)

	nextpi := findNextPlayer(me)
	if nextpi < 0 {
		return me.gotoPhase(_P_ROUNDS_SETTLEMENT)
	}

	me.unsetTimer(pd.timer)
	pd.timer = me.setTimer(30000, prtOnTimeOut)

	return me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_START_TURN,
		Source:  fwb.CMD_SOURCE_APP,
		Payload: me.turnOrder[pd.hotIndex],
	})
}

func prtOnTimeOut(me *gameImp, command sgs.Command) (bool, *er.Err) {
	fillTurn(me)
	return false, nil
}

func prtOnAction(me *gameImp, command sgs.Command) (bool, *er.Err) {
	pd := me.pd.(*prtData)

	if command.Source != me.turnOrder[pd.hotIndex] {
		return true, er.Throw(fwb.E_CMD_INVALID_CLIENT, er.EInfo{
			"details":        "Command source is not a valid client ID, or the client is not the currently enabled player",
			"ID":             command.Source,
			"current player": me.turnOrder[pd.hotIndex],
		}).To(me.lg)
	}

	action, err := me.ap.Parse(command)

	if err.Importance() >= er.IMPT_THREAT || action == nil {
		return true, err
	}

	if !action.ValidateAgainst(&me.gd) {
		return true, err.Push(me.app.SendToPlayer(command.Source, sgs.Command{
			ID:      fwb.CMD_ACTION_REJECTED,
			Source:  fwb.CMD_SOURCE_APP,
			Payload: command.Payload,
		}))
	}

	err = err.Push(action.Do(&me.gd))

	if err.Importance() >= er.IMPT_THREAT {
		return true, err.Push(me.app.SendToPlayer(command.Source, sgs.Command{
			ID:      fwb.CMD_ACTION_REJECTED,
			Source:  fwb.CMD_SOURCE_APP,
			Payload: command.Payload,
		}))
	}

	err = err.Push(me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_ACTION_COMMITTED,
		Source:  fwb.CMD_SOURCE_APP,
		Payload: me.gd,
	}))

	return true, err.Push(nextTurn(me))
}
