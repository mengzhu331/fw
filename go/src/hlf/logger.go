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
		me.conf = _defaultLogConf
	}
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
	text := me.formati(nil, lv, format, a...)
	if me.conf.ToConsole {
		me.send2console(text)
	}

	if me.conf.ToFile {
		if me.conf.Lv >= lv {
			text = me.formati(me, lv, format, a...)
			me.send2file(text, lv)
		}

		for log := me.parent; log != nil; log = log.parent {
			clv := me.findAppliedLv(log)
			if clv >= lv {
				text := me.formati(log, lv, format, a...)
				log.send2file(text, lv)
			}
		}
	}
}

func (me *logger) findAppliedLv(ancestor *logger) logLevel {
	var lv logLevel = LvUnknown
	found := false

	for log := me; log != ancestor && log.parent != nil; log = log.parent {
		lv, found = ancestor.conf.ChildLv[log.id]
		if found {
			break
		}
	}

	if !found {
		lv = ancestor.conf.DefaultChildLv
	}

	return lv
}

func (me *logger) send2console(text string) {
	me.send2Srv("console:", text)
}

func (me *logger) send2file(text string, lv logLevel) {
	me.send2Srv(me.getConsoleFileTarget(), text)
	if LvError >= lv {
		me.send2Srv(me.getErrorFileTarget(), text)
	}
}

func (me *logger) send2Srv(target string, text string) {
	li := logItem{
		target: target,
		text:   text,
	}
	_logSrvCh <- li
}

func (me *logger) toPrefix() string {
	if me.id == "" {
		return " "
	}
	return "[" + me.id + "] "
}

func (me *logger) formati(parent *logger, lv logLevel, format string, a ...interface{}) string {
	text := ""
	text += lv.toPrefix()
	text += logTime()
	text += me.indent(parent)
	text += me.toPrefix()
	text += fmt.Sprintf(format, a...)
	text += "\n"
	return text
}

func (me *logger) getIndent(parent *logger) int {
	indent := 0
	for log := me; log != parent && log.parent != nil; log = log.parent {
		indent += _logSysConf.Indent
	}
	return indent
}

func (me *logger) indent(parent *logger) string {
	indent := me.getIndent(parent)
	indents := ""
	for i := 0; i < indent; i++ {
		indents += " "
	}
	return indents
}

func (me *logger) getPath() string {
	path := ""
	for l := me; l != nil; l = l.parent {
		if l.id != "" {
			path = l.id + "/" + path
		}
	}
	path = _logRoot + path
	return path
}

func (me *logger) getConsoleFileTarget() string {
	return "file:" + me.getPath() + "console.log"
}

func (me *logger) getErrorFileTarget() string {
	return "file:" + me.getPath() + "error.log"
}

func logTime() string {
	return "[" + time.Now().Format(time.RFC3339) + "]"
}
