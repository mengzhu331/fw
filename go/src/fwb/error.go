package main

import "fmt"
import "runtime"

const (
	EMC_ERROR int = 0x80000000

	EMC_INTERACTION_ERROR int = EMC_ERROR | 0x1000000
	EMC_INTERNAL_ERROR    int = EMC_ERROR | 0x2000000

	EMC_ILLEGAL_PARAMETTER int = 0x1000
	EMC_STATE_INVALID      int = 0x2000

	//Interaction errors
	EC_COMMAND_NOT_FOR_THE_STATE int = EMC_INTERACTION_ERROR | EMC_ILLEGAL_PARAMETTER | 1
	EC_ILLEGAL_PLAYER_ID         int = EMC_INTERACTION_ERROR | EMC_ILLEGAL_PARAMETTER | 2
	EC_USE_FULL_AC               int = EMC_INTERACTION_ERROR | EMC_STATE_INVALID | 3
	EC_USE_AC_NOT_AFFORD         int = EMC_INTERACTION_ERROR | EMC_STATE_INVALID | 4

	//Internal errors
	EC_CURRENT_PHASE_INVALID int = EMC_INTERNAL_ERROR | EMC_STATE_INVALID | 1
	EC_GAME_CONTEXT_INVALID  int = EMC_INTERNAL_ERROR | EMC_STATE_INVALID | 2
	EC_ILLEGAL_LOCAL_STATE   int = EMC_INTERNAL_ERROR | EMC_STATE_INVALID | 3

	EC_NULL_PARAMETERS             int = EMC_INTERNAL_ERROR | EMC_ILLEGAL_PARAMETTER | 1
	EC_UNKNOWN_TARGET_PHASE        int = EMC_INTERNAL_ERROR | EMC_ILLEGAL_PARAMETTER | 2
	EC_ILLEGAL_PLAYER_NAME         int = EMC_INTERNAL_ERROR | EMC_ILLEGAL_PARAMETTER | 3
	EC_INVALID_APP_THIS            int = EMC_INTERNAL_ERROR | EMC_ILLEGAL_PARAMETTER | 4
	EC_ILLEGAL_COMMAND_TARGET      int = EMC_INTERNAL_ERROR | EMC_ILLEGAL_PARAMETTER | 5
	EC_ILLEGAL_COMMAND_PARTICIPANT int = EMC_INTERNAL_ERROR | EMC_ILLEGAL_PARAMETTER | 6
)

const (
	FWERROR_MAX_FRAMES int = 10
)

type FwError struct {
	ErrorCode int
	Message   string
	IntParam  int
	StrParam  string
	Nested    error
	Frames    *runtime.Frames
}

func (this *FwError) DefaultMessage() string {
	switch this.ErrorCode {
	case EC_COMMAND_NOT_FOR_THE_STATE:
		return fmt.Sprintf("Command received cannot be exectuted in the current state.(Command Id: %X, Command Source: %q)", this.IntParam, this.StrParam)
	case EC_UNKNOWN_TARGET_PHASE:
		return fmt.Sprintf("The name of the target phase is not recognized.(Target Phase Name: %q)", this.StrParam)
	default:
		return "No error message."
	}
}

func (this *FwError) Error() string {
	message := this.DefaultMessage()
	if this.Message != "" {
		message = this.Message
	}
	errorString := fmt.Sprintf("Error occurred with ErrorCode:0x%X\nError Message:%s\nnested error:",
		this.ErrorCode, message)
	if this.Nested != nil {
		errorString += "\n" + this.Nested.Error()
	} else {
		errorString += "nil\n"
	}

	return errorString
}

func MakeFwErrorComplete(code int, message string, intParam int, strParam string, nested error) *FwError {
	var frames *runtime.Frames

	pc := make([]uintptr, FWERROR_MAX_FRAMES)
	n := runtime.Callers(1, pc)
	if n > 0 {
		pc = pc[1:n]
		frames = runtime.CallersFrames(pc)
	}

	return &FwError{
		code,
		message,
		intParam,
		strParam,
		nested,
		frames,
	}
}

func MakeFwError(code int, intParam int, strParam string, nested error) *FwError {
	return MakeFwErrorComplete(code, "", intParam, strParam, nested)
}

func MakeFwErrorByCode(code int) *FwError {
	return MakeFwError(code, 0, "", nil)
}
