package log

import (
	"fmt"
	"runtime"
)

func getCaller() string {
	const skip = 2
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", f.Name(), line)
}
