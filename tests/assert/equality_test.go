package assert_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestNotEqualAndEqual(t *testing.T) {
	var pointer *http.Response
	tests := []struct {
		name     string
		expected any
		actual   any
		fails    bool
	}{
		{
			name:     "type-mismatch",
			expected: fmt.Errorf("%w", errors.New("smth")),
			actual:   "something",
			fails:    false,
		},
		{
			name:     "same-different-case",
			expected: "Something",
			actual:   "something",
			fails:    false,
		},
		{
			name:     "exactly-the-same",
			expected: "test",
			actual:   "test",
			fails:    true,
		},
		{
			name:     "interface-pointer-nil",
			expected: nil,
			actual:   pointer,
			fails:    false,
		},
	}

	for _, tt := range tests {
		t.Run("[NotEqual]: "+tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.NotEqual(test, tt.expected, tt.actual, "not equal expectation failed")

			assert.Equal(t, tt.fails, test.Failed())
		})
	}

	for _, tt := range tests {
		t.Run("[Equal]: "+tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.Equal(test, tt.expected, tt.actual, "not equal expectation failed")

			if tt.fails && test.Failed() {
				t.Error("failed expectation")
			}
		})
	}
}

func TestNotNil(t *testing.T) {
	tests := []struct {
		name  string
		v     any
		fails bool
	}{
		{
			name:  "error",
			v:     fmt.Errorf("%w", errors.New("smth")),
			fails: false,
		},
		{
			name:  "reflect-zero-value",
			v:     reflect.ValueOf(nil),
			fails: false,
		},
		{
			name:  "nil",
			v:     nil,
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.NotNil(test, tt.v)

			assert.Equal(t, tt.fails, test.Failed())
		})
	}
}

func TestNilError(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		fails bool
	}{
		{
			name:  "error",
			err:   fmt.Errorf("%w", errors.New("smth")),
			fails: true,
		},
		{
			name:  "nil",
			err:   nil,
			fails: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.NilError(test, tt.err)

			assert.Equal(t, tt.fails, test.Failed())
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		fails bool
	}{
		{
			name:  "error",
			err:   fmt.Errorf("%w", errors.New("smth")),
			fails: false,
		},
		{
			name:  "nil",
			err:   nil,
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.Error(test, tt.err)

			assert.Equal(t, tt.fails, test.Failed())
		})
	}
}

func TestErrorWith(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contents string
		fails    bool
	}{
		{
			name:     "without-substring",
			err:      fmt.Errorf("%w", errors.New("smth")),
			contents: "something",
			fails:    true,
		},
		{
			name:     "contains-substring",
			err:      errors.New("test"),
			contents: "test",
			fails:    false,
		},
		{
			name:     "nil",
			err:      nil,
			contents: "test",
			fails:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.ErrorWith(test, tt.err, tt.contents)

			assert.Equal(t, tt.fails, test.Failed())
		})
	}
}

func TestErrorIs(t *testing.T) {
	errInvalid := errors.New("invalid")

	tests := []struct {
		name   string
		err    error
		target error
		fails  bool
	}{
		{
			name:   "same",
			err:    errInvalid,
			target: errInvalid,
			fails:  false,
		},
		{
			name:   "not-exact-same",
			err:    errInvalid,
			target: errors.New("invalid"),
			fails:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.ErrorIs(test, tt.err, tt.target)

			assert.Equal(t, tt.fails, test.Failed())
		})
	}
}

func TestErrorType(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		fails  bool
	}{
		{
			name:   "same",
			err:    errors.New("invalid"),
			target: errors.New("invalid"),
			fails:  false,
		},
		{
			name: "same-type-other-props",
			err:  &json.UnmarshalTypeError{},
			target: &json.UnmarshalTypeError{
				Field: "f",
				Value: "bool",
			},
			fails: false,
		},
		{
			name: "different-type",
			err:  &json.InvalidUnmarshalError{},
			target: &json.UnmarshalTypeError{
				Field: "f",
				Value: "bool",
			},
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.ErrorType(test, tt.err, tt.target)

			assert.Equal(t, tt.fails, test.Failed())
		})
	}
}

func TestApproximately(t *testing.T) {
	tests := []struct {
		name     string
		expected time.Duration
		actual   time.Duration
		delta    time.Duration
		fails    bool
	}{
		{
			name:     "success-same",
			expected: time.Second,
			actual:   time.Second,
			delta:    time.Hour,
			fails:    false,
		},
		{
			name:     "success-approx-0.5s-delta",
			expected: time.Second,
			actual:   time.Second + (time.Second / 2),
			delta:    time.Second,
			fails:    false,
		},
		{
			name:     "success-approx-1s-delta",
			expected: time.Second,
			actual:   time.Second,
			delta:    time.Second,
			fails:    false,
		},
		{
			name:     "invalid-too-small-delta",
			expected: time.Second,
			actual:   time.Second + 2000,
			delta:    time.Microsecond,
			fails:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.Approximately(test, tt.expected, tt.actual, tt.delta)

			assert.Equal(t, tt.fails, test.Failed())
		})
	}
}
