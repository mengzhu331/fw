package main

import (
	"er"
	"hlf"
	"sgs"
	"strconv"
)

var _execMap = map[int]func(*fwAppImp, sgs.Command) *er.Err{
	sgs.CMD_TICK:    forwardToGame,
	sgs.CMD_APP_RUN: forwardToGame,
}

var _execTypeMap = map[int]func(*fwAppImp, sgs.Command) *er.Err{
	sgs.CMD_C_CLIENT: forwardToPlayer,
}

type fwApp interface {
	getGame() game
	getPlayers() map[int]player
	getSession() sgs.Session
	getLogger() hlf.Logger
}

type fwAppImp struct {
	gm game
	pm map[int]player
	s  sgs.Session
	lg hlf.Logger
}

func fwAppBuildFunc() sgs.App {
	return &fwAppImp{}
}

func (me *fwAppImp) Init(s sgs.Session, clients []int, profile string) *er.Err {
	me.s = s

	me.lg = s.GetLogger().Child("FamilyWarApp")

	me.pm = make(map[int]player)

	for c := range clients {
		me.pm[c] = makePlayer()
		me.pm[c].init(me, c)
	}

	me.gm = makeGame()
	return me.gm.init(me)
}

func (me *fwAppImp) SendCommand(command sgs.Command) *er.Err {
	exec, found := _execMap[command.ID]

	if found {
		return exec(me, command)
	}

	exec, found = _execTypeMap[command.ID&sgs.CMD_CATEGORY]

	if found {
		return exec(me, command)
	}

	return er.Throw(_E_INVALID_CMD, er.EInfo{
		"details": "command is invalid",
		"command": command.HexID(),
	}).To(me.lg)
}

func (me *fwAppImp) getGame() game {
	return me.gm
}

func (me *fwAppImp) getPlayers() map[int]player {
	return me.pm
}

func (me *fwAppImp) getSession() sgs.Session {
	return me.s
}

func (me *fwAppImp) getLogger() hlf.Logger {
	return me.lg
}

func forwardToGame(app *fwAppImp, command sgs.Command) *er.Err {
	return app.getGame().sendCommand(command)
}

func forwardToPlayer(app *fwAppImp, command sgs.Command) *er.Err {
	p, found := app.pm[command.Source]
	if !found {
		return er.Throw(_E_CMD_INVALID_CLIENT, er.EInfo{
			"details": "command from invalid client " + strconv.Itoa(command.Source) + ", command " + command.HexID(),
		}).To(app.lg)
	}

	return p.sendCommand(command)
}
