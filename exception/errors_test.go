package exception

import (
	"errors"
	"strings"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestAggregatef(t *testing.T) {
	tests := []struct {
		name    string
		errors  []error
		wantErr error
	}{
		{
			name:    "nil",
			wantErr: nil,
		},
		{
			name: "errors",
			errors: []error{
				errors.New("test"),
				errors.New("test2"),
			},
			wantErr: errors.New("[\"test\", \"test2\"]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Aggregatef(tt.errors...)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestStackTrace(t *testing.T) {
	got := StackTrace()

	lines := strings.Split(got, "\n")
	if len(lines) < 6 {
		t.Error("too few lines")
	}

	const expected = "exception.TestStackTrace"
	contains := false
	for _, v := range lines {
		if strings.Contains(v, expected) {
			contains = true
			break
		}
	}
	assert.Equal(t, true, contains, "stack trace does not contain current function path")
}
