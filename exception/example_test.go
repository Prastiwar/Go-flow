package exception_test

import (
	"errors"
	"fmt"

	"github.com/Prastiwar/Go-flow/exception"
)

func ExampleAggregatedError() {
	aggErr := exception.Aggregate(
		errors.New("'' is not valid value for title"),
		errors.New("title is required"),
	)
	err := fmt.Errorf("(%v) validation errors: %w", len(aggErr), aggErr)
	fmt.Println(err)

	// Output:
	// (2) validation errors: ["'' is not valid value for title", "title is required"]
}

func ExampleHandlePanicError() {
	defer exception.HandlePanicError(func(err error) {
		fmt.Println(err)
	})

	panic("string error")

	// Output:
	// string error
}
