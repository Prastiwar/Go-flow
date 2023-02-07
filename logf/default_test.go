package logf_test

import (
	"context"
	"testing"

	"github.com/Prastiwar/Go-flow/logf"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestDefault(t *testing.T) {
	t.Cleanup(func() {
		logf.SetDefault(nil)
	})

	l := logf.Default()
	assert.Equal(t, logf.NewLogger(), l)

	customLogger := logf.NewLogger(logf.WithFields(logf.Fields{"custom": "test"}))
	logf.SetDefault(func() logf.Logger {
		return customLogger
	})

	l = logf.Default()
	assert.Equal(t, customLogger, l)
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	l := logf.From(ctx)
	assert.Equal(t, logf.NewLogger(), l)

	customLogger := logf.NewLogger(logf.WithFields(logf.Fields{"custom": "test"}))
	ctx = logf.WithLogger(ctx, customLogger)
	assert.Equal(t, customLogger, logf.From(ctx))
}
