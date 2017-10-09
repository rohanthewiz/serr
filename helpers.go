package serr

import (
	"fmt"
	"runtime"
	"path/filepath"
)

func FunctionLoc() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}
