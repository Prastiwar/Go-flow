package assert

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

type ErrorFunc func(t *testing.T, err error)

type ResultErrorFunc[T any] func(t *testing.T, result T, err error)

func errorf(t *testing.T, msg string, prefixes ...string) {
	t.Helper()
	prefix := ""
	if len(prefixes) > 0 {
		prefix = strings.Join(prefixes, ": ") + ": "
	}
	t.Error(prefix + msg)
}

// Equal asserts expected and actual values are equal using deep equal from reflection.
func Equal(t *testing.T, expected interface{}, actual interface{}, prefixes ...string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		if expected == nil && reflect.ValueOf(actual).IsNil() {
			return
		}
		errorf(t, fmt.Sprintf("expected: '%v', actual: '%v'", expected, actual), prefixes...)
	}
}

// NotEqual asserts expected and actual values are not equal using equal operator.
func NotEqual(t *testing.T, expected interface{}, actual interface{}, prefixes ...string) {
	t.Helper()
	if expected == actual {
		errorf(t, fmt.Sprintf("expected: '%v', actual: '%v'", expected, actual), prefixes...)
	}
}

func NotNil(t *testing.T, val interface{}, prefixes ...string) {
	t.Helper()
	if val == nil {
		errorf(t, "expected not to be nil", prefixes...)
	}
}

func NilError(t *testing.T, err error, prefixes ...string) {
	t.Helper()
	if err != nil {
		errorf(t, fmt.Sprintf("expected error to be nil but was %v", err), prefixes...)
	}
}

func Error(t *testing.T, err error, prefixes ...string) {
	t.Helper()
	if err == nil {
		errorf(t, "expected error but got nil", prefixes...)
	}
}

// ErrorWith asserts err does contain target content within error string representation.
func ErrorWith(t *testing.T, err error, content string, prefixes ...string) {
	t.Helper()
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
	t.Helper()
	if !errors.Is(err, target) {
		errorf(t, fmt.Sprintf("expected '%#v' error but got '%#v'", target, err), prefixes...)
	}
}

// ErrorType matches just type of error.
func ErrorType(t *testing.T, err error, target error, prefixes ...string) {
	t.Helper()
	if reflect.TypeOf(err) != reflect.TypeOf(target) {
		errorf(t, fmt.Sprintf("expected '%#v' error but got '%#v'", target, err), prefixes...)
	}
}

// Approximately asserts actual duration is approximately(within delta difference) equal to expected duration.
func Approximately(t *testing.T, expected time.Duration, actual time.Duration, delta time.Duration, prefixes ...string) {
	t.Helper()
	if expected < actual-delta || expected > actual+delta {
		errorf(t, fmt.Sprintf("expected approximately '%v' duration but got '%v'", expected, actual), prefixes...)
	}
}
