package serr

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	Caller             = 1
	CallersParent      = 2
	CallersGrandParent = 3
)

// Return caller or ancestors calling location
// Optional level can be
// 1 - immediate caller (default)
// 2 - the callers parent
// 3 - the callers grandparent -- you get the idea
func FuncLoc(callerLevel ...int) string {
	lvl := Caller
	if len(callerLevel) > 0 {
		lvl = callerLevel[0]
	}

	_, file, line, ok := runtime.Caller(lvl)
	if !ok {
		return "could not determine location"
	}
	return PathLevel(fmt.Sprintf("%s:%d", file, line))
}

// Return a portion of the fullpath
// 0 - file only  // e.g. main.go
// 1 - file and parent  // e.g. myproject/main.go
// 2 - file up to grandparent  // e.g. githubusername/myproject/main.go
// defaults to 1 - parent/file
func PathLevel(path string, level ...uint) (subpath string) {
	lvl := 1
	if len(level) > 0 {
		lvl = int(level[0])
	}
	if path == "" {
		return path
	}

	tokens := strings.Split(path, string(filepath.Separator))
	ln := len(tokens)
	if ln <= 1 {
		fmt.Println("path not split")
		return path
	}

	idx := len(tokens) - int(lvl) - 1
	if idx < 0 {
		idx = 0
	}
	return strings.Join(tokens[idx:], string(filepath.Separator))
}

// This function is deprecated. Please use serr.FuncLoc()
func FunctionLoc() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}
