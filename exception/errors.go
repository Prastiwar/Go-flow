// Package exception provides helper functions to facilitate work with errors, stack traces or handling panics.
package exception

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Aggregate returns formatted array of errors as single error.
func Aggregate(errors ...error) error {
	count := len(errors)
	if count == 0 {
		return nil
	}

	b := strings.Builder{}

	b.WriteString("\"")
	b.WriteString(errors[0].Error())
	b.WriteString("\"")
	for i := 1; i < count; i++ {
		b.WriteString(", \"")
		b.WriteString(errors[i].Error())
		b.WriteString("\"")
	}

	return fmt.Errorf("[%v]", b.String())
}

// StackTrace returns a formatted string stack trace of the goroutine that calls it.
func StackTrace() string {
	return string(debug.Stack())
}
