package logf

import (
	"context"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestDefault(t *testing.T) {
	l := Default()
	assert.Equal(t, NewLogger(), l)

	customLogger := NewLogger(WithFields(Fields{"custom": "test"}))
	SetDefault(func() Logger {
		return customLogger
	})

	l = Default()
	assert.Equal(t, customLogger, l)
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	l := From(ctx)
	assert.Equal(t, NewLogger(), l)

	customLogger := NewLogger(WithFields(Fields{"custom": "test"}))
	ctx = WithLogger(ctx, customLogger)
	assert.Equal(t, customLogger, From(ctx))
}
