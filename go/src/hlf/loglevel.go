package hlf

type logLevel int

const (
	//_LV_UNKNOWN log level is unknown
	_LV_UNKNOWN = 0

	//_LV_ERROR log level for errors
	_LV_ERROR = iota * 10

	//_LV_WARNING log level for warnings
	_LV_WARNING

	//_LV_NOTIFICATION log level for notificaitons
	_LV_NOTIFICATION

	//_LV_INFO log level for information
	_LV_INFO

	//_LV_DEBUG log level for debug information
	_LV_DEBUG

	//LvTrace log level for trace information
	_LV_TRACE
)

func (me *logLevel) toPrefix() string {
	switch *me {
	case _LV_UNKNOWN:
		return "[Unknown]"
	case _LV_ERROR:
		return "[Error  ]"
	case _LV_WARNING:
		return "[Warning]"
	case _LV_NOTIFICATION:
		return "[Note   ]"
	case _LV_INFO:
		return "[Info   ]"
	case _LV_DEBUG:
		return "[Debug  ]"
	case _LV_TRACE:
		return "[Trace  ]"
	}
	return "[Unknown]"
}
