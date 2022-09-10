package assert

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type ErrorFunc func(t *testing.T, err error)

type ResultErrorFunc[T any] func(t *testing.T, result T, err error)

func errorf(t *testing.T, msg string, prefixes ...string) {
	prefix := ""
	if len(prefixes) > 0 {
		prefix = strings.Join(prefixes, ": ") + ": "
	}
	t.Error(prefix + msg)
}

// Equal asserts expected and actual values are equal using deep equal from reflection
func Equal(t *testing.T, expected interface{}, actual interface{}, prefixes ...string) {
	if !reflect.DeepEqual(expected, actual) {
		errorf(t, fmt.Sprintf("expected: '%v', actual: '%v'", expected, actual), prefixes...)
	}
}

// Equal asserts expected and actual values are not equal using equal operator
func NotEqual(t *testing.T, expected interface{}, actual interface{}, prefixes ...string) {
	if expected == actual {
		errorf(t, fmt.Sprintf("expected: '%v', actual: '%v'", expected, actual), prefixes...)
	}
}

func NotNil(t *testing.T, val interface{}, prefixes ...string) {
	if val == nil {
		errorf(t, "expected not to be nil", prefixes...)
	}
}

func NilError(t *testing.T, err error, prefixes ...string) {
	if err != nil {
		errorf(t, fmt.Sprintf("expected error to be nil but was %v", err), prefixes...)
	}
}

func Error(t *testing.T, err error, prefixes ...string) {
	if err == nil {
		errorf(t, "expected error but got nil", prefixes...)
	}
}

// ErrorWith asserts err does contain target content within error string representation.
func ErrorWith(t *testing.T, err error, content string, prefixes ...string) {
	errFormat := "expected error with content: '%v', error: '%v'"
	if err == nil {
		errorf(t, fmt.Sprintf(errFormat, content, err), prefixes...)
		return
	}

	ok := strings.Contains(err.Error(), content)
	if !ok {
		errorf(t, fmt.Sprintf(errFormat, content, err.Error()), prefixes...)
	}
}

// ErrorIs asserts whether error in err's chain matches target.
func ErrorIs(t *testing.T, err error, target error, prefixes ...string) {
	if !errors.Is(err, target) {
		errorf(t, fmt.Sprintf("expected '%#v' error but got '%#v'", target, err), prefixes...)
	}
}

// ErrorType matches just type of error.
func ErrorType(t *testing.T, err error, target error, prefixes ...string) {
	if reflect.TypeOf(err) != reflect.TypeOf(target) {
		errorf(t, fmt.Sprintf("expected '%#v' error but got '%#v'", target, err), prefixes...)
	}
}
