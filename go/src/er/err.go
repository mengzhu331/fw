package er

import (
	"encoding/json"
	"fmt"
	"hlf"
	"runtime"
	"strconv"
)

//Err error data
type Err struct {
	code       int32
	callStack  []string
	stackDepth int
	info       string
	next       *Err
}

func (me *Err) Error() string {
	if me == nil {
		return "No Error"
	}

	return fmt.Sprintf("Error Code: 0x%x, Error Info: %v", me.code, me.info)
}

//Code retrieve code
func (me *Err) Code() int32 {
	if me == nil {
		return 0
	}

	return me.code
}

//Importance retrieve error importance
func (me *Err) Importance() int32 {
	if me == nil {
		return IMPT_NONE
	}

	return me.code & E_IMPORTANCE
}

//DumpCallStack generate a call stack integrated string
func (me *Err) DumpCallStack(lv int) string {
	if me == nil {
		return ""
	}

	if lv > len(me.callStack) {
		lv = len(me.callStack)
	}
	ds := "Call Stack:\n"
	callstack := me.callStack[:lv]

	for i := 0; i < len(callstack); i++ {
		ds += callstack[i]
		if i != len(callstack)-1 {
			ds += " by:\n"
		}
	}

	if len(callstack) < me.stackDepth {
		ds += "by:\n"
		ds += "  ..." + strconv.Itoa(me.stackDepth-len(callstack)) + " more"
	}

	return ds
}

//EInfo map for info entries
type EInfo map[string]interface{}

var _maxCallStackFrames = 100

//Throw init an Err
func Throw(code int32, info EInfo) *Err {
	errinfo, _ := json.Marshal(info)

	pc := make([]uintptr, _maxCallStackFrames)
	n := runtime.Callers(1, pc)
	var frames *runtime.Frames
	if n > 0 {
		pc = pc[1:n]
		frames = runtime.CallersFrames(pc)
	}

	callstack := make([]string, 0)

	for {
		frame, next := frames.Next()
		framefootprint := fmt.Sprintf("  [%v(), line %v] called, in [%v]", frame.Function, frame.Line, frame.File)
		callstack = append(callstack, framefootprint)
		if !next {
			break
		}
	}

	return &Err{
		code:       code,
		callStack:  callstack,
		stackDepth: n - 1,
		info:       string(errinfo),
	}
}

//To log error
func (me *Err) To(logger hlf.Logger) *Err {

	if me == nil {
		return me
	}

	if (me.code | E_IMPORTANCE) >= IMPT_THREAT {
		logger.Err(me.Error())
		logger.To("error.log").Err(me.Error() + ", " + me.DumpCallStack(10))
		return me
	}

	if (me.code | E_IMPORTANCE) >= IMPT_REMARKABLE {
		logger.Ntf(me.Error())
		logger.To("remarkable.log").Ntf(me.Error() + ", " + me.DumpCallStack(10))
		return me
	}

	logger.Inf(me.Error())
	logger.To("exception.log").Inf(me.Error() + ", " + me.DumpCallStack(10))
	return me
}

//Push support errors stack, add one error to the top
func (me *Err) Push(top *Err) *Err {
	if top == nil {
		return me
	}

	prev := &Err{
		next: me,
	}

	for prev.next != nil && prev.next.Importance() > top.Importance() {
		prev = prev.next
	}

	top.next = prev.next
	prev.next = top

	return prev.next
}

//Pop support errors stack, remove the top error
func (me *Err) Pop() *Err {
	return me.next
}
