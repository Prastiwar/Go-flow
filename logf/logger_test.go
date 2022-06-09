package logf

import (
	"goflow/tests/assert"
	"goflow/tests/mocks"
	"log"
	"testing"
)

func TestWithScope(t *testing.T) {
	loggerMock := log.New(nil, "", 0)
	fields := Fields{"1": "1"}

	logger := WithScope(loggerMock, fields)
	writer := logger.Writer().(formatterWriter)
	assert.MapMatch(t, writer.fields, fields)

	assert.NotNil(t, logger)
}

func TestWithFormatter(t *testing.T) {
	writerCounter := assert.Count(1)
	formatterCounter := assert.Count(1)

	writerMock := mocks.NewWriterMock(func(p []byte) (n int, err error) {
		writerCounter.Inc()
		return 0, nil
	})
	loggerMock := log.New(writerMock, "", 0)
	formatter := NewFormatterMock(func(msg string, fields Fields) string {
		formatterCounter.Inc()
		return msg
	})

	logger := WithFormatter(loggerMock, formatter)
	Info(logger, "test")

	assert.NotNil(t, logger)
	formatterCounter.Assert(t)
	writerCounter.Assert(t)
}

func TestInfo(t *testing.T) {
	logTest(t, Info)
}

func TestInfof(t *testing.T) {
	logTestf(t, Infof)
}

func TestWarn(t *testing.T) {
	logTest(t, Warn)
}

func TestWarnf(t *testing.T) {
	logTestf(t, Warnf)
}

func TestError(t *testing.T) {
	logTest(t, Error)
}

func TestErrorf(t *testing.T) {
	logTestf(t, Errorf)
}

func TestDebug(t *testing.T) {
	logTest(t, Debug)
}

func TestDebugf(t *testing.T) {
	logTestf(t, Debugf)
}

func TestTrace(t *testing.T) {
	logTest(t, Trace)
}

func TestTracef(t *testing.T) {
	logTestf(t, Tracef)
}

func TestFatal(t *testing.T) {
	resetLogger()
	counter := assert.Count(1)
	writerMock := mocks.NewWriterMock(func(p []byte) (n int, err error) {
		counter.Inc()
		t.Skip()
		return 0, nil
	})
	loggerMock := log.New(writerMock, "", 0)

	Fatal(loggerMock, "test")

	counter.Assert(t)
}

func TestFatalf(t *testing.T) {
	resetLogger()
	counter := assert.Count(1)
	writerMock := mocks.NewWriterMock(func(p []byte) (n int, err error) {
		counter.Inc()
		t.Skip()
		return 0, nil
	})
	loggerMock := log.New(writerMock, "", 0)

	Fatalf(loggerMock, "%v", "test")

	counter.Assert(t)
}

func logTest(t *testing.T, fn func(*log.Logger, interface{})) {
	resetLogger()
	counter := assert.Count(1)
	writerMock := mocks.NewWriterMock(func(p []byte) (n int, err error) {
		counter.Inc()
		return 0, nil
	})
	loggerMock := log.New(writerMock, "", 0)

	fn(loggerMock, "test")

	counter.Assert(t)
}

func logTestf(t *testing.T, fn func(*log.Logger, string, ...any)) {
	resetLogger()
	counter := assert.Count(1)
	writerMock := mocks.NewWriterMock(func(p []byte) (n int, err error) {
		counter.Inc()
		return 0, nil
	})
	loggerMock := log.New(writerMock, "", 0)

	fn(loggerMock, "%v", "test")

	counter.Assert(t)
}
