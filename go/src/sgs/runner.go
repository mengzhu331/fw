package sgs

import (
	"hlf"
)

var _log hlf.Logger = hlf.MakeLogger("SGS")
var _srv *sessionServer

//Run turn on the sgs servers
func Run(abf AppBuildFunc, confPath string) error {

	_log.Inf("Starting SGS...")

	if len(confPath) >= 1 && confPath[len(confPath)-1] != '/' && confPath[len(confPath)-1] != '\\' {
		confPath += "/"
	}

	c := loadConf(confPath + "sgs.conf")

	var e error
	_srv, e = makeSSrv(SSrvParam{
		Profile:        c.App.Profile,
		DefaultClients: c.App.DefaultClients,
		MinimalClients: c.App.MinimalClients,
		OptimalWS:      c.App.OptimalWaitSecond,
		BaseTickMs:     c.BaseTickMs,
		ABF:            abf,
	})

	if e != nil {
		_log.Err("Failed to init SSVR, SGS shut down")
		return e
	}

	sgasPrx := createSgasPrx()
	sgasPrx.setServerURI(c.AuthSrvURI)

	e = webStartUp(WebSrvParam{
		Port:        c.Port,
		WSReadBuff:  c.WSReadBuff,
		WSWriteBuff: c.WSWriteBuff,
	}, sgasPrx)

	if e != nil {
		_log.Err("Failed to start SGS Web Server, SGS shut down")
	} else {
		_log.Inf("SGS started")
	}

	return e
}
