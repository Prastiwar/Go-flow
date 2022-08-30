package exception_test

import (
	"errors"
	"fmt"

	"github.com/Prastiwar/Go-flow/exception"
)

func ExampleAggregate() {
	aggErr := exception.Aggregate(
		errors.New("'' is not valid value for title"),
		errors.New("title is required"),
	)
	err := fmt.Errorf("validation: %w", aggErr)
	fmt.Println(err)

	// Output:
	// validation: ["'' is not valid value for title", "title is required"]
}

func ExampleHandlePanicError() {
	defer exception.HandlePanicError(func(err error) {
		fmt.Println(err)
	})

	panic("string error")

	// Output:
	// string error
}
