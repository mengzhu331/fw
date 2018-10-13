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
		hotIndex: 0,
		timer:    -1,
	}

	me.pd = pd

	return nextTurn(me)
}

func switchHotIndex(me *gameImp) {
	pd := me.pd.(*prtData)
	pd.hotIndex++
	if pd.hotIndex >= len(me.app.GetPlayers()) {
		pd.hotIndex = 0
	}
}

func nextTurn(me *gameImp) *er.Err {
	pd := me.pd.(*prtData)

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

	switchHotIndex(me)
	return true, err.Push(nextTurn(me))
}
