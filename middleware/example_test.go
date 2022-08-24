package middleware_test

import (
	"errors"
	"fmt"
	"goflow/middleware"
)

func Example() {
	type pipeRequest string
	type pipeResponse error

	// middleware pipeline for request of type 'pipeRequest' and response of type 'pipeResponse'
	middleware := middleware.NewMiddleware[pipeRequest, pipeResponse]()

	middleware.Use(
		func(r pipeRequest, next func(r pipeRequest) pipeResponse) pipeResponse {
			fmt.Println("1")
			response := next(r)
			fmt.Println("4")
			return response
		},
	)

	middleware.Use(
		func(r pipeRequest, next func(r pipeRequest) pipeResponse) pipeResponse {
			fmt.Println("2")
			validate := func(r pipeRequest) bool { return true }

			ok := validate(r)
			if !ok {
				// stop pipeline and return error
				return errors.New("validation failed")
			}
			return next(r)
		},
	)

	handler := func(r pipeRequest) pipeResponse {
		fmt.Println("3")
		return nil
	}

	// wrap middleware to handler
	wrappedHandler := middleware.Wrap(handler)
	request := pipeRequest("request")

	// run pipeline
	response := wrappedHandler(request)
	fmt.Println("Response:", response)

	// Output:
	// 1
	// 2
	// 3
	// 4
	// Response: <nil>
}
