package main

import (
	"hlf"
	"sgs"
)

var _log hlf.Logger = hlf.CreateLogger("FWB")

var _ch = make(chan string)

func main() {
	_log.Inf("FWB starting up")
	e := sgs.Run(fwAppBuildFunc)

	if e == nil {
		var c string

		for c != "quit" {
			select {
			case c = <-_ch:
			}
		}
	}
	_log.Inf("FWB shut down")
}
