package main

import (
	"er"
	"hlf"
	"sgs/ssvr"
)

type phaseDataMap map[int]interface{}

type game interface {
	init(app fwApp)
	sendCommand(command ssvr.Command) *er.Err
}

type gameImp struct {
	app fwApp
	gd  gameData
	pd  phaseDataMap
	phs phase
	lg  hlf.Logger
}

type execCmd func(*gameImp, ssvr.Command) *er.Err

type enterPhase func(*gameImp)

var _globalCmdMap = map[int]execCmd{
	ssvr.CMD_APP_RUN: run,
}

var _phaseEnterMap = map[phase]enterPhase{}

var _phaseCmdMap = map[phase]map[int]execCmd{}

var _defaultCmdMap = map[int]execCmd{
	ssvr.CMD_TICK: onTickDefault,
}

func (me *gameImp) init(app fwApp) {
	me.app = app
	me.pd = make(phaseDataMap)
	me.gd.round = 0

	me.gd.pData = make(map[int]playerData)

	for k := range app.getPlayers() {
		pdata := initPlayerData()
		pdata.playerID = k
		me.gd.pData[k] = pdata
	}

	me.gd.cards = make([]cardData, 0)

	me.lg = app.getLogger()
}

func (me *gameImp) sendCommand(command ssvr.Command) *er.Err {

	exec, found := _globalCmdMap[command.ID]

	if found {
		return exec(me, command)
	}

	pcm, pfound := _phaseCmdMap[me.phs]

	if pfound {
		exec, found = pcm[command.ID]
		if found {
			return exec(me, command)
		}
	}

	exec, found = _defaultCmdMap[command.ID]

	if found {
		return exec(me, command)
	}

	return er.Throw(_E_CMD_NOT_EXEC, er.EInfo{
		"details": "gameImp is not supposed to receive the command",
		"command": ssvr.CmdHexID(command),
		"phase":   me.phs,
	}).To(me.lg)
}

func (me *gameImp) gotoPhase(p phase) {
	me.phs = p
	exec, found := _phaseEnterMap[p]
	if found {
		exec(me)
	}
}

func run(me *gameImp, command ssvr.Command) *er.Err {
	me.gotoPhase(_P_GAME_START)
	return nil
}

func makeGame() game {
	return &gameImp{}
}

func onTickDefault(me *gameImp, command ssvr.Command) *er.Err {
	return nil
}
