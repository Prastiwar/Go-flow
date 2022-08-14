package exception

import (
	"errors"
	"fmt"
)

func HandlePanicError(onPanic func(error)) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				onPanic(v)
			case string:
				onPanic(errors.New(v))
			default:
				onPanic(errors.New(fmt.Sprint(v)))
			}
		}
	}()
}
