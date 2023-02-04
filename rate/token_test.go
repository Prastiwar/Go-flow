package rate_test

import (
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

var _ rate.Token = &falseToken{}

type falseToken struct{}

func (*falseToken) Allow() bool {
	return false
}

func (*falseToken) ResetsAt() time.Time {
	return rate.MaxTime
}

func (*falseToken) Use() error {
	return rate.ErrInvalidTokenValue
}

func TestFalseToken(t *testing.T) {
	assert.NotNil(t, rate.FalseToken)
	assert.Equal(t, false, rate.FalseToken.Allow())
	assert.Equal(t, rate.ErrInvalidTokenValue, rate.FalseToken.Use())
	assert.Equal(t, rate.MaxTime, rate.FalseToken.ResetsAt())
}

func TestIsFalseToken(t *testing.T) {
	tests := []struct {
		name  string
		token rate.Token
		want  bool
	}{
		{
			name:  "true",
			token: rate.FalseToken,
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
			got := rate.IsFalseToken(tt.token)
			assert.Equal(t, tt.want, got)
		})
	}
}
