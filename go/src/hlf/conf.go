package hlf

type lvMap map[string]logLevel

type confMap map[string]loggerConf

//conf framework settings
type conf struct {
	LogRoot     string
	Indent      int
	DefaultFile string
	Loggers     confMap
}

//loggerConf logger settings
type loggerConf struct {
	ToFile         bool
	ToConsole      bool
	Lv             logLevel
	DefaultChildLv logLevel
	ChildLv        lvMap
}

var _defaultLogConf = loggerConf{
	ToFile:         true,
	ToConsole:      true,
	Lv:             _LV_INFO,
	DefaultChildLv: _LV_NOTIFICATION,
	ChildLv:        make(lvMap),
}

var _conf = conf{
	LogRoot:     "./log/",
	Indent:      2,
	DefaultFile: "console.log",

	Loggers: make(confMap),
}
