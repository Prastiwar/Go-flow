package middleware_test

import (
	"errors"
	"testing"

	"github.com/Prastiwar/Go-flow/middleware"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

type pipeRequest string

type pipeResponse error

func TestPipelineOrder(t *testing.T) {
	expectedError := errors.New("end-error")
	expectedMessages := []string{"log-start", "pre-handler", "handler", "post-handler", "log-end"}
	actualMessages := make([]string, 0)

	middleware := middleware.NewMiddleware[pipeRequest, pipeResponse]()

	middleware.Use(
		func(r pipeRequest, next func(r pipeRequest) pipeResponse) pipeResponse {
			actualMessages = append(actualMessages, "log-start")
			err := next(r)
			actualMessages = append(actualMessages, "log-end")
			return err
		},
	)

	middleware.Use(
		func(r pipeRequest, next func(r pipeRequest) pipeResponse) pipeResponse {
			actualMessages = append(actualMessages, "pre-handler")
			return next(r)
		},
	)

	middleware.Use(
		func(r pipeRequest, next func(r pipeRequest) pipeResponse) pipeResponse {
			_ = next(r)
			actualMessages = append(actualMessages, "post-handler")
			return expectedError
		},
	)

	handler := func(r pipeRequest) pipeResponse {
		actualMessages = append(actualMessages, "handler")
		return nil
	}

	wrappedHandler := middleware.Wrap(handler)
	err := wrappedHandler("")

	assert.Equal(t, expectedError, err)
	assert.Equal(t, len(expectedMessages), len(actualMessages))
	for i, v := range expectedMessages {
		assert.Equal(t, v, actualMessages[i])
	}
}
