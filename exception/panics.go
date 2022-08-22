package exception

import (
	"errors"
	"fmt"
)

// HandlePanicError defers function which check for recover() and convert it to error preserving
// error type or creating new string error. If recover() is nil onPanic will not be called.
func HandlePanicError(onPanic func(error)) {
	if r := recover(); r != nil {
		onPanic(ConvertToError(r))
	}
}

// ConvertToError returns i if it's error or new error string.
func ConvertToError(i any) error {
	switch v := i.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return errors.New(fmt.Sprint(v))
	}
}
