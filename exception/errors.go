// Package exception provides helper functions to facilitate work with errors, stack traces or handling panics.
package exception

import (
	"fmt"
	"strings"
)

// Aggregate returns formatted array of errors as single error.
func Aggregate(errors ...error) error {
	count := len(errors)
	if count == 0 {
		return nil
	}

	b := strings.Builder{}

	b.WriteString(errors[0].Error())
	for i := 1; i < count; i++ {
		b.WriteString(", ")
		b.WriteString(errors[i].Error())
	}

	return fmt.Errorf("[%v]", b.String())
}
