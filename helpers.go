package serr

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

var FrameLevels = frameLevel{0, 1, 2, 3, 4, 5}

// frameLevel is a representation of the current stack frame level
// as a distance from the frame getting some kind of context from the runtime
type frameLevel struct {
	// Same frame
	FrameLevel0,
	// Parent
	FrameLevel1,
	// Grandparent
	FrameLevel2,
	// great grandparent
	FrameLevel3,
	// great great grandparent
	FrameLevel4,
	// great great great grandparent
	FrameLevel5 int
}

// StringFromErr returns an enriched string if err is a SErr,
// or the standard error string otherwise
func StringFromErr(err error) (strErr string) {
	if err == nil {
		return
	}
	strErr = err.Error()
	if ser, ok := err.(SErr); ok {
		strErr = ser.String()
	}
	return
}

// UserMsgFromErr returns the user message in the SErr,
// alt string if none, empty string if no error
func UserMsgFromErr(err error, alt ...string) (msg string) {
	if err == nil {
		return
	}

	if ser, ok := err.(SErr); ok {
		msg, _ = ser.UserMsg()
	}

	if msg == "" && len(alt) > 0 {
		msg = alt[0]
	}
	return
}

// FunctionLoc returns last two path tokens of caller.
// optFuncLevel passes the function level to go back up.
// The default is 1, referring to the caller of this function
func FunctionLoc(optFuncLevel ...int) string {
	frameLevel := 1 // default to the caller's frame
	if len(optFuncLevel) > 0 {
		frameLevel = optFuncLevel[0]
	}

	_, fPath, line, ok := runtime.Caller(frameLevel)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s:%d", LastNTokens(fPath, "/", 2), line)
}

// FunctionName returns the function name of the caller
// optFuncLevel passes the function level to go back up.
// The default is 1, referring to the caller of this function
func FunctionName(optFuncLevel ...int) (funcName string) {
	frameLevel := 1 // default to the caller's frame
	if len(optFuncLevel) > 0 {
		frameLevel = optFuncLevel[0]
	}

	if pc, _, _, ok := runtime.Caller(frameLevel); ok {
		fPtr := runtime.FuncForPC(pc)
		if fPtr == nil {
			return
		}

		return LastNTokens(fPtr.Name(), "/", 2)
	}
	return
}

func LastNTokens(str, separator string, n int) (lastTokens string) {
	tokens := strings.Split(str, "/")

	if len(tokens) >= n {
		lastTokens = strings.Join(tokens[len(tokens)-n:], separator)
	} else {
		lastTokens = filepath.Base(str)
	}
	return
}
