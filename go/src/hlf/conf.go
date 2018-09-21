package hlf

import (
	"encoding/json"
	"io/ioutil"
)

type lvMap map[string]logLevel

type confMap map[string]LoggerConf

//Conf framework settings
type Conf struct {
	LogRoot     string
	Indent      int
	DefaultFile string

	Loggers confMap
}

//LoggerConf logger settings
type LoggerConf struct {
	ToFile         bool
	ToConsole      bool
	Lv             logLevel
	DefaultChildLv logLevel
	ChildLv        lvMap
}

var _defaultLogConf = LoggerConf{
	ToFile:         true,
	ToConsole:      true,
	Lv:             LvInfo,
	DefaultChildLv: LvNotification,
	ChildLv:        make(lvMap),
}

var _conf = Conf{
	LogRoot:     "./log/",
	Indent:      2,
	DefaultFile: "console.log",

	Loggers: make(confMap),
}

func loadConfFile(f string, conf interface{}) error {
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(c, conf)
	return err
}
