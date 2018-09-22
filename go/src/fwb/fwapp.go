package main

import (
	"er"
	"sgs/ssvr"
)

type fwApp struct{}

func fwAppBuildFunc() ssvr.App {
	return &fwApp{}
}

func (me *fwApp) Init(s ssvr.Session, clients []ssvr.NetClient) *er.Err {
	return nil
}

func (me *fwApp) SendCommand(command ssvr.Command) *er.Err {
	return nil
}
