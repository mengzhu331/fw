package core

import (
	"er"
	"fwb"
	"sgs"
)

type prtData struct {
	hotIndex int
	timer    int
	turn     int
}

func prtInit(me *gameImp) *er.Err {
	me.lg.Dbg("Enter Round Turns phase")

	pd := &prtData{
		hotIndex: len(me.app.GetPlayers()) - 1,
		timer:    -1,
		turn:     0,
	}

	me.pd = pd

	me.setDCE(fwb.CMD_ACTION, prtOnAction)

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

	pd.turn++
	me.alg.Inf("Turn %v", pd.turn)

	if pd.timer >= 0 {
		me.unsetTimer(pd.timer)
	}
	pd.timer = me.setTimer(2000, prtOnTimeOut)

	return me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_START_TURN,
		Who:     fwb.CMD_WHO_APP,
		Payload: me.turnOrder[pd.hotIndex],
	})
}

func prtOnTimeOut(me *gameImp, command sgs.Command) *er.Err {
	pd := me.pd.(*prtData)
	me.app.SendToMockPlayer(me.turnOrder[pd.hotIndex], sgs.Command{
		ID:      fwb.CMD_START_TURN,
		Who:     fwb.CMD_WHO_APP,
		Payload: me.turnOrder[pd.hotIndex],
	})
	return nil
}

func prtOnAction(me *gameImp, command sgs.Command) *er.Err {
	pd := me.pd.(*prtData)

	if command.Who != me.turnOrder[pd.hotIndex] {
		return er.Throw(fwb.E_CMD_INVALID_CLIENT, er.EInfo{
			"details":        "Command source is not a valid client ID, or the client is not the currently enabled player",
			"ID":             command.Who,
			"current player": me.turnOrder[pd.hotIndex],
		}).To(me.lg)
	}

	action, err := me.ap.Parse(command)

	if err.Importance() >= er.IMPT_THREAT || action == nil {
		return err
	}

	if !action.ValidateAgainst(&me.gd) {
		return err.Push(me.app.SendToPlayer(command.Who, sgs.Command{
			ID:      fwb.CMD_ACTION_REJECTED,
			Who:     fwb.CMD_WHO_APP,
			Payload: command.Payload,
		}))
	}

	err = err.Push(action.Do(&me.gd))

	printAction(me, action)

	if err.Importance() >= er.IMPT_THREAT {
		return err.Push(me.app.SendToPlayer(command.Who, sgs.Command{
			ID:      fwb.CMD_ACTION_REJECTED,
			Who:     fwb.CMD_WHO_APP,
			Payload: command.Payload,
		}))
	}

	printTurnInfo(me)

	err = err.Push(me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_ACTION_COMMITTED,
		Who:     command.ID,
		Payload: me.gd,
	}))

	return err.Push(nextTurn(me))
}

func printAction(me *gameImp, action fwb.Action) {
	me.alg.Inf(action.String())
}

func printTurnInfo(me *gameImp) {
	me.alg.Inf("Player Data")
	for _, p := range me.gd.PData {
		me.alg.Inf("  Player %v: %v", me.app.GetPlayer(p[fwb.PD_CLIENT_ID]).Name(), p[fwb.PD_CLIENT_ID+1:])
	}
}
