package main

import (
	"err"
	"hlf"
	"time"
)

const (
	_ETest = err.ImptRecoverable | err.ETInternal | 0x1
)

func main() {

	log0 := hlf.CreateLogger("")
	log1 := hlf.CreateLogger("sys")
	log2 := log1.Child("module")
	log3 := log2.Child("function")

	e := err.Throw(_ETest, err.EInfo{"param": "ok"})
	e.To(log3)

	log0.Err("default logger error")
	log1.Err("sys logger error")
	log2.Err("module logger error")
	log3.Err("function logger error")

	log0.Wrn("default logger warning")
	log1.Wrn("sys logger warning")
	log2.Wrn("module logger warning")
	log3.Wrn("function logger warning")

	log0.Ntf("default logger notification")
	log1.Ntf("sys logger notification")
	log2.Ntf("module logger notification")
	log3.Ntf("function logger notification")

	log0.Inf("default logger info")
	log1.Inf("sys logger info")
	log2.Inf("module logger info")
	log3.Inf("function logger info")

	log0.Dbg("default logger debug")
	log1.Dbg("sys logger debug")
	log2.Dbg("module logger debug")
	log3.Dbg("function logger debug")

	log0.Trc("default logger trace")
	log1.Trc("sys logger trace")
	log2.Trc("module logger trace")
	log3.Trc("function logger trace")

	<-time.After(time.Duration(2) * time.Second)
	/*
		s := "{\"ID\": 1,\"Param\":{\"t\": \"s\"}}"

		type Cmd struct {
			ID    int
			Param interface{}
		}

		c := Cmd{}

		json.Unmarshal([]byte(s), &c)

		sgs.Run(fwAppBuildFunc)
	*/
}
