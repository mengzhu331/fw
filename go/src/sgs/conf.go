package sgs

import (
	"er"
	"hlf"
)

//webConf web server settings
type webConf struct {
	Port        int
	EPM         map[string]string
	WSReadBuff  int
	WSWriteBuff int
}

//sessionConf game session settings
type sessionConf struct {
	DefaultClientsPerSession int
	BaseTickMs               int
}

//conf sgs web server configuration
type conf struct {
	Web     webConf
	Session sessionConf
}

var _defaultConf = conf{
	Web: webConf{
		Port:        9090,
		EPM:         nil,
		WSReadBuff:  1024,
		WSWriteBuff: 1024,
	},

	Session: sessionConf{
		DefaultClientsPerSession: 3,
		BaseTickMs:               100,
	},
}

var _defaultEPM = map[string]string{
	"login":        "/login",
	"ws":           "/ws",
	"join_session": "/join_game",
}

//loadConf read conf file to get the settings, default values will be used when a field is missing
func loadConf(f string) conf {
	c := _defaultConf

	e := hlf.LoadConfFile(f, &c)

	if e != nil {
		er.Throw(_E_LOAD_CONF_FAIL, er.EInfo{"filename": f}).To(_log)
	}

	if c.Web.EPM == nil {
		c.Web.EPM = make(map[string]string)
	}

	for k, v := range _defaultEPM {
		_, found := c.Web.EPM[k]
		if !found {
			c.Web.EPM[k] = v
		}
	}
	return c
}
