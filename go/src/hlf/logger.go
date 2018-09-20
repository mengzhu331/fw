package hlf

import (
	"fmt"
	"time"
)

//Logger log mechanism interface
type Logger interface {
	Ntf(string, ...interface{})
	Inf(string, ...interface{})
	Err(string, ...interface{})
	Wrn(string, ...interface{})
	Dbg(string, ...interface{})
	Trc(string, ...interface{})
}

func (me *logger) loadConf() {
	id := me.id
	if id == "" {
		id = "_"
	}
	conf, found := _conf[id]
	if found {
		me.conf = conf
	} else if me.parent != nil {
		me.conf = me.parent.conf
	} else {
		me.conf = _defaultConf
	}
}

//CreateLogger init a logger
func CreateLogger(id string, parent Logger) Logger {
	lg := logger{
		id: id,
	}
	var islogger bool
	lg.parent, islogger = parent.(*logger)
	if !islogger {
		lg.parent = nil
	}
	lg.loadConf()
	return &lg
}

//CreateRootLogger init a logger
func CreateRootLogger(id string) Logger {
	return CreateLogger(id, nil)
}

//CreateDefaultLogger init a logger
func CreateDefaultLogger() Logger {
	return CreateLogger("", nil)
}

type logger struct {
	id     string
	conf   Conf
	parent *logger
}

func (me *logger) Ntf(format string, a ...interface{}) {
	me.print(LvNotification, format, a...)
}

func (me *logger) Inf(format string, a ...interface{}) {
	me.print(LvInfo, format, a...)
}

func (me *logger) Err(format string, a ...interface{}) {
	me.print(LvError, format, a...)
}

func (me *logger) Wrn(format string, a ...interface{}) {
	me.print(LvWarning, format, a...)
}

func (me *logger) Dbg(format string, a ...interface{}) {
	me.print(LvDebug, format, a...)
}

func (me *logger) Trc(format string, a ...interface{}) {
	me.print(LvTrace, format, a...)
}

func (me *logger) print(lv logLevel, format string, a ...interface{}) {
	text := formati(nil, me, lv, format, a...)
	if me.conf.ToConsole {
		sendToSrv("console:", text)
	}

	if me.conf.ToFile {
		if me.conf.Lv >= lv {
			sendToSrv(getConsoleFileTarget(me), text)
			if LvError >= lv {
				sendToSrv(getErrorFileTarget(me), text)
			}
		}

		log := me.parent
		nc := me
		for log != nil {
			clv, found := log.conf.ChildLv[nc.id]
			if !found {
				clv = log.conf.DefaultChildLv
			}

			if clv >= lv {
				text := formati(log, me, lv, format, a...)
				sendToSrv(getConsoleFileTarget(log), text)
				if LvError >= lv {
					sendToSrv(getErrorFileTarget(log), text)
				}
			}

			log = log.parent
			nc = log
		}
	}
}

func sendToSrv(target string, text string) {
	li := logItem{
		target: target,
		text:   text,
	}
	_ch <- li
}

func (me *logger) toPrefix() string {
	if me.id == "" {
		return ""
	}
	return "[" + me.id + "]"
}

func formati(parent *logger, log *logger, lv logLevel, format string, a ...interface{}) string {
	text := ""
	text += lv.toPrefix()
	text += logTime()
	text += log.toPrefix()
	text += indent(parent, log)
	text += fmt.Sprintf(format, a...)
	text += "\n"
	return text
}

func getIndent(parent *logger, log *logger) int {
	indent := 0
	for log != parent && log != nil {
		indent += _logSysConf.Indent
		log = log.parent
	}
	return indent
}

func indent(parent *logger, log *logger) string {
	indent := getIndent(parent, log)
	indents := ""
	for i := 0; i < indent; i++ {
		indents += " "
	}
	return indents
}

func logTime() string {
	return "[" + time.Now().Format(time.RFC3339) + "]"
}

func getPath(log *logger) string {
	path := ""
	for l := log; l != nil; l = l.parent {
		if l.id != "" {
			path = l.id + "/" + path
		}
	}
	path = _logRoot + path
	return path
}

func getConsoleFileTarget(log *logger) string {
	return "file:" + getPath(log) + "console.log"
}

func getErrorFileTarget(log *logger) string {
	return "file:" + getPath(log) + "error.log"
}
