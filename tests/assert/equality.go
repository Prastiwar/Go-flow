package assert

import (
	"reflect"
	"strings"
	"testing"
)

type ErrorFunc func(t *testing.T, err error)

// Equal asserts expected and actual values are equal using deep equal from reflection
func Equal(t *testing.T, expected interface{}, actual interface{}, prefixes ...string) {
	if !reflect.DeepEqual(expected, actual) {
		prefix := ""
		if len(prefixes) > 0 {
			prefix = strings.Join(prefixes, ": ") + ": "
		}
		t.Errorf("%vexpected: '%v', actual: '%v'", prefix, expected, actual)
	}
}

// Equal asserts expected and actual values are not equal using equal operator
func NotEqual(t *testing.T, expected interface{}, actual interface{}, prefixes ...string) {
	if expected == actual {
		prefix := ""
		if len(prefixes) > 0 {
			prefix = strings.Join(prefixes, ": ") + ": "
		}
		t.Errorf("%vexpected: '%v', actual: '%v'", prefix, expected, actual)
	}
}

func NotNil(t *testing.T, val interface{}) {
	if val == nil {
		t.Error("expected not nil")
	}
}

func NilError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func Error(t *testing.T, err error) {
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func ErrorWith(t *testing.T, err error, content string) {
	errFormat := "expected error with content: '%v', error: '%v'"
	if err == nil {
		t.Errorf(errFormat, content, err)
		return
	}

	ok := strings.Contains(err.Error(), content)
	if !ok {
		t.Errorf(errFormat, content, err.Error())
	}
}
