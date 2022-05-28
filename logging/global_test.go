package logging

import (
	"errors"
	"goflow/tests/assert"
	"testing"
)

func TestWithOptions(t *testing.T) {
	WithOptions(
		WithLogFormat("test"),
	)

	assert.Equal(t, "test", defaultLog.Options().logf)
}

func TestInfo(t *testing.T) {
	logger, callCounter := mockLog(t, "test")
	defaultLog = logger

	Info("test")

	callCounter.Assert(t)
}

func TestWarn(t *testing.T) {
	logger, callCounter := mockLog(t, "test")
	defaultLog = logger

	Warn("test")

	callCounter.Assert(t)
}

func TestError(t *testing.T) {
	err := errors.New("test")
	logger, callCounter := mockLog(t, err)
	defaultLog = logger

	Error(err)

	callCounter.Assert(t)
}
