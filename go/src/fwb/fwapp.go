package main

import (
	"sgs/ssvr"
)

type fwApp struct{}

func fwAppBuildFunc() ssvr.App {
	return &fwApp{}
}

func (me *fwApp) Init(c chan string, clients []ssvr.NetClient) error {
	return nil
}

func (me *fwApp) SendCommand(command ssvr.Command) error {
	return nil
}
