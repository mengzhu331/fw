package hlf

import (
	"encoding/json"
	"io/ioutil"
)

type lvMap map[string]logLevel

type confMap map[string]Conf

//Conf logger settings
type Conf struct {
	ToFile         bool
	ToConsole      bool
	Lv             logLevel
	DefaultChildLv logLevel
	ChildLv        lvMap
}

var _defaultConf = Conf{
	ToFile:         true,
	ToConsole:      true,
	Lv:             LvInfo,
	DefaultChildLv: LvNotification,
	ChildLv:        make(lvMap),
}

var _conf = make(confMap)

//LogSysConf log system settings
type LogSysConf struct {
	LogRoot string
	Indent  int
}

var _logSysConf = LogSysConf{
	LogRoot: "./log/",
	Indent:  2,
}

func loadLogConf(fsys string, f string) bool {
	loadError := loadConfFile(fsys, &_logSysConf) != nil
	loadError = (loadConfFile(f, &_conf) != nil) || loadError
	return !loadError
}

func loadConfFile(f string, conf interface{}) error {
	c, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(c, conf)
	if err != nil {
		panic(err.Error())
	}
	return err
}
