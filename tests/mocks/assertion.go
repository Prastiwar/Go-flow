package mocks

import (
	"fmt"
	"runtime"
	"strings"
)

// unexpectedCall panics with unexpected call message. Error return value is for API simplicity.
func unexpectedCall(prefixes ...string) error {
	prefix := ""
	if len(prefixes) > 0 {
		prefix = strings.Join(prefixes, ": ") + ": "
	}

	msg := prefix + fmt.Sprintf("unexpected call on '%v'", callerName())
	panic(msg)
}

func callerName() string {
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		n := details.Name()
		start := strings.LastIndex(n, "/") + 1
		return n[start:]
	}
	return ""
}
