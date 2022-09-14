// Package exception provides helper functions to facilitate work with errors, e traces or handling panics.
package exception

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// AggregatedError is an array of errors. It implements error interface and aggregates multiple errors.
type AggregatedError []error

// Aggregate returns AggregatedError filled with specified errors. This is helper function
// to use variadic errors as cast to AggregatedError.
func Aggregate(errs ...error) AggregatedError {
	return errs
}

// Flat returns AggregatedError with flat one-dimensional array. Any nested AggregatedError will be unwrapped.
func (err AggregatedError) Flat() AggregatedError {
	for i := 0; i < len(err); i++ {
		current := err[i]

		if agg, ok := current.(AggregatedError); ok {
			err = delete(err, i)
			err = insert(err, i, agg)
			continue
		}
	}

	return err
}

func (err AggregatedError) Error() string {
	return Aggregatedf(err...).Error()
}

// Delete removes the element at i index from s, returning the modified slice.
func delete[S ~[]E, E any](s S, i int) S {
	return append(s[:i], s[i+1:]...)
}

// insert is copied function from https://cs.opensource.google/go/x/exp/+/dc92f865:slices/slices.go;l=136
func insert[S ~[]E, E any](s S, i int, v []E) S {
	tot := len(s) + len(v)
	if tot <= cap(s) {
		s2 := s[:tot]
		copy(s2[i+len(v):], s[i:])
		copy(s2[i:], v)
		return s2
	}
	s2 := make(S, tot)
	copy(s2, s[:i])
	copy(s2[i:], v)
	copy(s2[i+len(v):], s[i:])
	return s2
}

// Aggregatedf returns formatted array of errors as single error.
func Aggregatedf(errors ...error) error {
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

// StackTrace returns a formatted string e trace of the goroutine that calls it.
func StackTrace() string {
	return string(debug.Stack())
}
