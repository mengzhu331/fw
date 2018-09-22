package main

import (
	"er"
	"sgs/ssvr"
)

type player interface {
	init(app fwApp, id int)
	sendCommand(command ssvr.Command) *er.Err
}

type playerImp struct {
	app fwApp
	id  int
}

func (me *playerImp) init(app fwApp, id int) {
	me.app = app
	me.id = id
}

func (me *playerImp) sendCommand(command ssvr.Command) *er.Err {
	return nil
}

func makePlayer() player {
	return &playerImp{}
}
