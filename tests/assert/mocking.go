package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// ExpectCall panics with unexpected call message for caller name if fn is nil. If fn is not
// a func then it will panic with invalid parameter message. Standard way to call this is:
//
//	func (m MyMock) Do(key string) error {
//		assert.ExpectCall(m.OnDo) // if m.OnDo is nil - it will panic with accurate message
//		return m.OnDo(key)
//	}
func ExpectCall(fn any) {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		panic("expectCall accepts only func kind value as parameter")
	}
	if v.IsNil() {
		unexpectedCallWith(prevCallerName())
	}
}

// unexpectedCallWith panics with unexpected call message with callerName.
func unexpectedCallWith(callerName string, prefixes ...string) {
	prefix := ""
	if len(prefixes) > 0 {
		prefix = strings.Join(prefixes, ": ") + ": "
	}

	msg := prefix + fmt.Sprintf("unexpected call on '%v'", callerName)
	panic(msg)
}

// prevCallerName returns a parent caller name for caller.
func prevCallerName() string {
	return callerName(2)
}

// callerName returns caller name with with skip option.
func callerName(skip int) string {
	pc, _, _, ok := runtime.Caller(1 + skip)
	details := runtime.FuncForPC(pc)
	name := ""
	start := 0
	if ok && details != nil {
		name = details.Name()
		start = strings.LastIndex(name, "/") + 1
	}
	return name[start:]
}
