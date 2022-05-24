package assert

import (
	"strings"
	"testing"
)

// Equal asserts expected and actual values are equal using equal operator
func Equal(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("expected: '%v', actual: '%v'", expected, actual)
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
	errFormat := "content: '%v', error: '%v'"
	if err == nil {
		t.Errorf(errFormat, content, err)
		return
	}

	ok := strings.Contains(err.Error(), content)
	if !ok {
		t.Errorf(errFormat, content, err.Error())
	}
}
