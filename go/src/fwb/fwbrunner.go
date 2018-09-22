package main

import (
	"hlf"
	"sgs"
)

var _log hlf.Logger = hlf.MakeLogger("FWB")

var _ch = make(chan string)

func main() {
	_log.Ntf("FWB starting up...")
	e := sgs.Run(fwAppBuildFunc)
	if e == nil {
		var c string

		_log.Ntf("FWB started, available for client connections")

		for c != "quit" {
			select {
			case c = <-_ch:
			}
		}
	}
	_log.Inf("FWB shut down")
}
