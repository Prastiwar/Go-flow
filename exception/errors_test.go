package exception

import (
	"errors"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestAggregate(t *testing.T) {
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
