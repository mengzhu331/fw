package sgs

import (
	"er"
	"sutil"
)

type appConf struct {
	Profile           string
	DefaultClients    int
	MinimalClients    int
	OptimalWaitSecond int
}

//conf sgs web server configuration
type conf struct {
	Port        int
	WSReadBuff  int
	WSWriteBuff int
	BaseTickMs  int
	AuthSrvURI  string
	App         appConf
}

var _defaultConf = conf{
	Port:        9090,
	WSReadBuff:  1024,
	WSWriteBuff: 1024,
	BaseTickMs:  100,
	AuthSrvURI:  "http://127.0.0.1:3115",
	App: appConf{
		Profile:           "2pvp",
		DefaultClients:    2,
		MinimalClients:    2,
		OptimalWaitSecond: 30,
	},
}

//loadConf read conf file to get the settings, default values will be used when a field is missing
func loadConf(f string) conf {
	c := _defaultConf

	e := sutil.LoadConfFile(f, &c)

	if e != nil {
		er.Throw(_E_LOAD_CONF_FAIL, er.EInfo{"filename": f, "system error": e.Error()}).To(_log)
	}

	return c
}
