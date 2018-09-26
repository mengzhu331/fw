package sgs

import (
	"hlf"
)

var _log hlf.Logger = hlf.MakeLogger("SGS")

//Run turn on the sgs servers
func Run(abf AppBuildFunc) error {

	_log.Inf("Starting SGS...")

	c := loadConf("./sgs.conf")

	e := initSSrv(SSrvParam{
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

	e = webStartUp(WebSrvParam{
		Port:        c.Port,
		WSReadBuff:  c.WSReadBuff,
		WSWriteBuff: c.WSWriteBuff,
	}, createSgasPrx())

	if e != nil {
		_log.Err("Failed to start SGS Web Server, SGS shut down")
	} else {
		_log.Inf("SGS started")
	}

	return e
}
