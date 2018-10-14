package core

import (
	"er"
	"fwb"
	"fwb/actn"
	"math"
	"sgs"
)

type mpCommand struct {
	clientID int
	command  sgs.Command
}

type mockPlayers struct {
	g            *gameImp
	commandQueue []mpCommand
}

type mpCmdExe func(*gameImp, int, sgs.Command) *er.Err

var _mpCmdExeMap = map[int]mpCmdExe{
	fwb.CMD_GAME_START: mpOnGameStart,

	fwb.CMD_GAME_OVER: mpOnGameOver,

	fwb.CMD_ROUND_START: mpOnRoundStart,

	fwb.CMD_START_TURN: mpOnStartTurn,

	fwb.CMD_ACTION_COMMITTED: mpOnActionCommitted,

	fwb.CMD_ROUND_SETTLEMENT: mpOnRoundsSettlement,

	fwb.CMD_ROUND_SETTLEMENT_UPDATE: mpOnRoundSettlementUpdate,

	fwb.CMD_GAME_FINISH: mpOnGameFinish,
}

func (me *mockPlayers) sendCommand(clientID int, command sgs.Command) {
	me.commandQueue = append(me.commandQueue, mpCommand{
		clientID: clientID,
		command:  command,
	})
}

func (me *mockPlayers) init(game *gameImp) {
	me.g = game
}

func (me *mockPlayers) execOne() *er.Err {
	if len(me.commandQueue) < 1 {
		return nil
	}

	cinf := me.commandQueue[0]
	me.commandQueue = me.commandQueue[1:]

	ce, found := _mpCmdExeMap[cinf.command.ID]

	if !found {
		return er.Throw(fwb.E_CMD_NOT_EXECUTABLE, er.EInfo{
			"details":   "mock players cannot execute the command",
			"commandID": cinf.command.ID,
		}).To(me.g.lg)
	}

	return ce(me.g, cinf.clientID, cinf.command)
}

func (me *mockPlayers) execQueue() *er.Err {
	var err *er.Err
	for len(me.commandQueue) > 0 {
		err = err.Push(me.execOne())
		if err.Importance() >= er.IMPT_DEGRADE {
			return err
		}
	}
	return err
}

func mpOnGameStart(me *gameImp, cid int, command sgs.Command) *er.Err {
	return me.SendCommand(sgs.Command{
		ID:     fwb.CMD_GAME_START_ACK,
		Source: cid,
	})
}

func mpOnGameOver(me *gameImp, cid int, command sgs.Command) *er.Err {
	return nil
}

func mpOnRoundStart(me *gameImp, cid int, command sgs.Command) *er.Err {
	return me.SendCommand(sgs.Command{
		ID:     fwb.CMD_ROUND_START_ACK,
		Source: cid,
	})
}

func mpOnStartTurn(me *gameImp, cid int, command sgs.Command) *er.Err {
	cp := command.Payload.(int)
	if cid == cp {
		type actnCmd struct {
			ActionID int
			Payload  interface{}
		}

		return me.SendCommand(sgs.Command{
			ID:     fwb.CMD_ACTION,
			Source: cid,
			Payload: actnCmd{
				ActionID: actn.ACTN_SKIP,
				Payload:  nil,
			},
		})
	}
	return nil
}

func mpOnActionCommitted(me *gameImp, cid int, command sgs.Command) *er.Err {
	return nil
}

func mpOnRoundsSettlement(me *gameImp, cid int, command sgs.Command) *er.Err {
	type sd struct {
		Cereals int
		Meat    int
		Sweater int
	}

	pl := sd{}

	pl.Meat = int(math.Min(5.0, float64(me.gd.PData[cid][fwb.PD_PT_MEAT])))
	pl.Cereals = int(math.Min(5.0-float64(pl.Meat), float64(me.gd.PData[cid][fwb.PD_PT_CEREALS])))
	pl.Sweater = int(math.Min(5.0, float64(me.gd.PData[cid][fwb.PD_PT_SWEATER])))

	return me.SendCommand(sgs.Command{
		ID:      fwb.CMD_COMMIT_ROUND_SETTLEMENT,
		Source:  cid,
		Payload: pl,
	})
}

func mpOnRoundSettlementUpdate(me *gameImp, cid int, command sgs.Command) *er.Err {
	return nil
}

func mpOnGameFinish(me *gameImp, cid int, command sgs.Command) *er.Err {
	return nil
}
