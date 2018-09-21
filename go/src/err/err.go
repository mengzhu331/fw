package err

import (
	"encoding/json"
	"fmt"
	"hlf"
	"runtime"
)

//Err error data
type Err struct {
	code      int32
	callStack []string
	info      string
}

func (me *Err) Error() string {
	return fmt.Sprintf("Error Code: 0x%x, Error Info: %v", me.code, me.info)
}

//DumpCallStack generate a call stack integrated string
func (me *Err) DumpCallStack(lv int) string {
	if lv > len(me.callStack) {
		lv = len(me.callStack)
	}
	ds := "Call Stack:\n"
	callstack := me.callStack[:lv]

	for i := 0; i < len(callstack); i++ {
		ds += callstack[i] + "\n"
	}
	return ds
}

//EInfo map for info entries
type EInfo map[string]interface{}

var _maxCallStackFrames = 10

//Throw init an Err
func Throw(code int32, info EInfo) Err {
	errinfo, _ := json.Marshal(info)

	pc := make([]uintptr, _maxCallStackFrames)
	n := runtime.Callers(1, pc)
	var frames *runtime.Frames
	if n > 0 {
		pc = pc[1:n]
		frames = runtime.CallersFrames(pc)
	}

	callstack := make([]string, 0)

	for frame, next := frames.Next(); next; frame, next = frames.Next() {
		framefootprint := fmt.Sprintf("  [%v(), line %v] called, in [%v], by:", frame.Function, frame.Line, frame.File)
		callstack = append(callstack, framefootprint)
	}

	return Err{
		code:      code,
		callStack: callstack,
		info:      string(errinfo),
	}
}

//To log error
func (me *Err) To(logger hlf.Logger) {
	logger.Err(me.Error())
	logger.To("error.log").Err(me.Error() + ", " + me.DumpCallStack(10))
}
