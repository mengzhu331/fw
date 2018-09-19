package sgs

import (
	"sgs/ssvr"
	"sgs/web"
)

//Run turn on the sgs servers
func Run(abf ssvr.AppBuildFunc) {
	c := LoadConf("./conf.json")
	ssvr.Init(ssvr.SSrvParam{
		CPS:        c.Session.CPS,
		BaseTickMs: c.Session.BaseTickMs,
		ABF:        abf,
	})

	web.StartUp(web.WebSrvParam{
		Port:        c.Web.Port,
		EPM:         c.Web.EPM,
		WSReadBuff:  c.Web.WSReadBuff,
		WSWriteBuff: c.Web.WSWriteBuff,
	})
}
