package logging

import (
	"errors"
	"goflow/tests/assert"
	"testing"
)

func mockLog(t *testing.T, expectedArgs ...any) (*LogMock, *assert.Counter) {
	callCounter := assert.Count(1)
	logFunc := func(level Level, format string, args ...any) {
		argsLen := len(expectedArgs)
		assert.Equal(t, argsLen, len(args))

		for i := 0; i < argsLen; i++ {
			assert.Equal(t, expectedArgs[i], args[i])
		}

		callCounter.Inc()
	}

	return NewLogMock(logFunc), callCounter
}

func TestLogger(t *testing.T) {
	logger := NewLogger()

	logger.Log(Errorl, "")
	opts := logger.Options()

	assert.Equal(t, DefaultOptions().logf, opts.logf)
}

func TestLogInfo(t *testing.T) {
	logger, callCounter := mockLog(t, "test")

	LogInfo(logger, "test")

	callCounter.Assert(t)
}

func TestLogWarn(t *testing.T) {
	logger, callCounter := mockLog(t, 1)

	LogWarn(logger, 1)

	callCounter.Assert(t)
}

func TestLogErr(t *testing.T) {
	err := errors.New("error")
	logger, callCounter := mockLog(t, err, 2)

	LogErr(logger, err, 2)

	callCounter.Assert(t)
}

func TestFormat(t *testing.T) {
	logger, _ := mockLog(t)

	provided := format(logger, "test")
	def := format(logger, "")

	assert.Equal(t, "test", provided)
	assert.Equal(t, logger.Options().logf, def)
}
