package main

import (
	"er"
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
}

func (me *gameImp) init(app fwApp) {
	me.app = app
	me.pd = make(phaseDataMap)
}

func (me *gameImp) sendCommand(command ssvr.Command) *er.Err {
	return nil
}

func makeGame() game {
	return &gameImp{}
}
