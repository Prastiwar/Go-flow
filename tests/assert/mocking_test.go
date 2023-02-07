package assert

import (
	"testing"
)

func TestExpectCall(t *testing.T) {
	tests := []struct {
		name        string
		fn          any
		expectation any
	}{
		{
			name:        "success-call",
			fn:          func() {},
			expectation: nil,
		},
		{
			name:        "failure-nil-func",
			fn:          (func())(nil),
			expectation: "unexpected call on 'assert.TestExpectCall.func2'",
		},
		{
			name:        "failure-nil-inalid-param",
			fn:          (*bool)(nil),
			expectation: "expectCall accepts only func kind value as parameter",
		},
		{
			name:        "failure-nil",
			fn:          nil,
			expectation: "expectCall accepts only func kind value as parameter",
		},
		{
			name:        "failure-invalid-param",
			fn:          "test",
			expectation: "expectCall accepts only func kind value as parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				Equal(t, tt.expectation, recover())
			}()
			ExpectCall(tt.fn)
		})
	}
}

func TestUnexpectedCallWith(t *testing.T) {
	defer func() {
		Equal(t, "test: unexpected call on 'caller'", recover())
	}()
	unexpectedCallWith("caller", "test")
}

func TestPrevCallerName(t *testing.T) {
	v := CallPrevCallerName()
	Equal(t, "assert.TestPrevCallerName", v)

	v = NestedCallPrevCallerName()
	Equal(t, "assert.NestedCallPrevCallerName", v)
}

func CallPrevCallerName() string {
	return prevCallerName()
}

func NestedCallPrevCallerName() string {
	return CallPrevCallerName()
}
