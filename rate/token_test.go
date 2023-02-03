package rate

import (
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestFalseToken(t *testing.T) {
	assert.NotNil(t, FalseToken)
	assert.Equal(t, false, FalseToken.Allow())
	assert.Equal(t, ErrInvalidTokenValue, FalseToken.Use())
	assert.Equal(t, MaxTime, FalseToken.ResetsAt())
}

func TestIsFalseToken(t *testing.T) {
	tests := []struct {
		name  string
		token Token
		want  bool
	}{
		{
			name:  "true",
			token: FalseToken,
			want:  true,
		},
		{
			name:  "false",
			token: &falseToken{},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFalseToken(tt.token)
			assert.Equal(t, tt.want, got)
		})
	}
}
