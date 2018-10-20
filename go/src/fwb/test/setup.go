package test

import (
	"hlf"
	"sutil"
)

type conf struct {
	Port int
}

var _conf conf

var _tl = hlf.MakeLogger("Test")

func loadConf() {
	sutil.LoadConfFile("sgs.conf", &_conf)
}
