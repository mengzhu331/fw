package sgs

import (
	"hlf"
	"sgs/ssvr"
	"sgs/web"
)

var _log hlf.Logger = hlf.MakeLogger("SGS Runner")

//Run turn on the sgs servers
func Run(abf ssvr.AppBuildFunc) error {

	_log.Inf("Starting SGS servers...")

	c := loadConf("./sgs.conf")

	e := ssvr.Init(ssvr.SSrvParam{
		Profile:        c.App.Profile,
		DefaultClients: c.App.DefaultClients,
		MinimalClients: c.App.MinimalClients,
		OptimalWS:      c.App.OptimalWaitSecond,
		BaseTickMs:     c.Session.BaseTickMs,
		ABF:            abf,
	})

	if e != nil {
		_log.Err("Failed to start SSVR, SGS shut down")
		return e
	}

	e = web.StartUp(web.WebSrvParam{
		Port:        c.Web.Port,
		EPM:         c.Web.EPM,
		WSReadBuff:  c.Web.WSReadBuff,
		WSWriteBuff: c.Web.WSWriteBuff,
	})

	if e != nil {
		_log.Err("Failed to start SGS Web Server, SGS shut down")
	} else {
		_log.Inf("SGS servers started")
	}

	return e
}
