package sgs

import (
	"sgs/web"

	"github.com/tkanos/gonfig"
)

//WebConf web server settings
type WebConf struct {
	Port        int
	EPM         web.EndPointMap
	WSReadBuff  int
	WSWriteBuff int
}

//SessionConf game session settings
type SessionConf struct {
	CPS        int
	BaseTickMs int
}

//Conf sgs web server configuration
type Conf struct {
	Web     WebConf
	Session SessionConf
}

var _defaultConf = Conf{
	Web: WebConf{
		Port:        9090,
		EPM:         nil,
		WSReadBuff:  1024,
		WSWriteBuff: 1024,
	},

	Session: SessionConf{
		CPS:        3,
		BaseTickMs: 100,
	},
}

var _defaultEPM = web.EndPointMap{
	web.EP_LOGIN:       "/login",
	web.EP_WEBSOCKET:   "/ws",
	web.EP_JOINSESSION: "/join_game",
}

//LoadConf read conf file to get the settings, default values will be used when a field is missing
func LoadConf(f string) Conf {
	c := _defaultConf

	gonfig.GetConf(f, &c)

	if c.Web.EPM == nil {
		c.Web.EPM = make(map[int]string)
	}

	for i := web.EP_MIN + 1; i < web.EP_MAX; i++ {
		_, found := c.Web.EPM[i]
		if !found {
			c.Web.EPM[i] = _defaultEPM[i]
		}
	}
	return c
}
