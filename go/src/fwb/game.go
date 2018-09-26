package main

import (
	"er"
	"hlf"
	"sgs"
	"strconv"
)

type phaseDataMap map[int]interface{}

type game interface {
	init(app fwApp) *er.Err
	sendCommand(command sgs.Command) *er.Err
}

type gameImp struct {
	app     fwApp
	gd      gameData
	pd      phaseDataMap
	phs     phase
	lg      hlf.Logger
	cm      *cardManager
	profile string
}

type execCmd func(*gameImp, sgs.Command) *er.Err

type enterPhase func(*gameImp)

var _globalCmdMap = map[int]execCmd{
	sgs.CMD_APP_RUN: run,
}

var _phaseEnterMap = map[phase]enterPhase{}

var _phaseCmdMap = map[phase]map[int]execCmd{}

var _defaultCmdMap = map[int]execCmd{
	sgs.CMD_TICK: onTickDefault,
}

func (me *gameImp) init(app fwApp) *er.Err {
	me.app = app
	me.pd = make(phaseDataMap)
	me.gd.round = 0

	me.gd.pData = make(map[int]playerData)

	for k := range app.getPlayers() {
		pdata := initPlayerData()
		pdata.playerID = k
		me.gd.pData[k] = pdata
	}

	me.lg = app.getLogger()

	if len(app.getPlayers()) == 3 {
		me.profile = _PROFILE_3PVP
	} else if len(app.getPlayers()) == 2 {
		me.profile = _PROFILE_2PVP
	} else {
		er.Throw(_E_NO_PROPER_PROFILE, er.EInfo{
			"details": "no proper profile for the game parameter",
			"player":  strconv.Itoa(len(app.getPlayers())),
		}).To(me.lg)
	}

	me.gd.cards = make([]card, 0)

	me.cm = &cardManager{
		gm: me,
	}

	e := me.cm.loadCards()

	return e
}

func (me *gameImp) sendCommand(command sgs.Command) *er.Err {

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
		"command": command.HexID(),
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

func run(me *gameImp, command sgs.Command) *er.Err {
	me.gotoPhase(_P_GAME_START)
	return nil
}

func makeGame() game {
	return &gameImp{}
}

func onTickDefault(me *gameImp, command sgs.Command) *er.Err {
	return nil
}
