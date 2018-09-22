package main

import (
	"er"
	"hlf"
	"sgs/ssvr"
	"strconv"
)

var _execMap = map[int]func(*fwAppImp, ssvr.Command) *er.Err{
	ssvr.CMD_TICK: forwardToGame,
}

var _execTypeMap = map[int]func(*fwAppImp, ssvr.Command) *er.Err{
	ssvr.CMD_C_CLIENT: forwardToPlayer,
}

type fwApp interface {
	getGame() game
	getPlayers() map[int]player
	getSession() ssvr.Session
	getLogger() hlf.Logger
}

type fwAppImp struct {
	gm game
	pm map[int]player
	s  ssvr.Session
	lg hlf.Logger
}

func fwAppBuildFunc() ssvr.App {
	return &fwAppImp{}
}

func (me *fwAppImp) Init(s ssvr.Session, clients []int) *er.Err {
	me.s = s

	me.lg = s.GetLogger().Child("FamilyWarApp")

	for c := range clients {
		me.pm[c] = makePlayer()
		me.pm[c].init(me, c)
	}

	me.gm = makeGame()
	me.gm.init(me)

	return nil
}

func (me *fwAppImp) SendCommand(command ssvr.Command) *er.Err {
	exec, found := _execMap[command.ID]

	if found {
		return exec(me, command)
	}

	exec, found = _execTypeMap[command.ID&ssvr.CMD_CATEGORY]

	if found {
		return exec(me, command)
	}

	return er.Throw(_E_INVALID_CMD, er.EInfo{
		"details": "command is invalid",
		"ID":      command.ID,
	}).To(me.lg)
}

func (me *fwAppImp) getGame() game {
	return me.gm
}

func (me *fwAppImp) getPlayers() map[int]player {
	return me.pm
}

func (me *fwAppImp) getSession() ssvr.Session {
	return me.s
}

func (me *fwAppImp) getLogger() hlf.Logger {
	return me.lg
}

func forwardToGame(app *fwAppImp, command ssvr.Command) *er.Err {
	return app.getGame().sendCommand(command)
}

func forwardToPlayer(app *fwAppImp, command ssvr.Command) *er.Err {
	p, found := app.pm[command.Source]
	if !found {
		return er.Throw(_E_CMD_INVALID_CLIENT, er.EInfo{
			"details": "command from invalid client " + strconv.Itoa(command.Source) + ", command " + strconv.Itoa(command.ID),
		}).To(app.lg)
	}

	return p.sendCommand(command)
}
