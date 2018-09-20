package hlf

type logLevel int

const (
	//LvUnknown log level is unknown
	LvUnknown = 0

	//LvError log level for errors
	LvError = iota * 10

	//LvWarning log level for warnings
	LvWarning

	//LvNotification log level for notificaitons
	LvNotification

	//LvInfo log level for information
	LvInfo

	//LvDebug log level for debug information
	LvDebug

	//LvTrace log level for trace information
	LvTrace
)

func (me *logLevel) toPrefix() string {
	switch *me {
	case LvUnknown:
		return "[Unknown]"
	case LvError:
		return "[Error  ]"
	case LvWarning:
		return "[Warning]"
	case LvNotification:
		return "[Note   ]"
	case LvInfo:
		return "[Info   ]"
	case LvDebug:
		return "[Debug  ]"
	case LvTrace:
		return "[Trace  ]"
	}
	return "[Unknown]"
}
