package serr

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Public convenience var for providing caller indirection typically for FunctionLoc
var CallerIndirection = struct {
	Caller, CallersParent, CallersGrandParent, CallersGreatGrandParent int
}{1, 2, 3, 4}

func FunctionLoc(indirection ...int) string {
	skip := 2
	if len(indirection) > 0 {
		skip = indirection[0]
	}

	const pathSeparator = "/"
	_, fPath, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	var partialPath string
	tokens := strings.Split(fPath, "/")
	if len(tokens) >= 2 {
		partialPath = strings.Join(tokens[len(tokens)-2:], pathSeparator)
	} else {
		partialPath = filepath.Base(fPath)
	}
	return fmt.Sprintf("%s:%d", partialPath, line)
}
