package logf

import (
	"goflow/tests/assert"
	"log"
	"testing"
	"time"
)

func TestSetFormatter(t *testing.T) {
	resetLogger()
	counter := assert.Count(1)
	formatter := NewFormatterMock(func(msg string, fields Fields) string {
		counter.Inc()
		return msg
	})

	SetFormatter(formatter)
	SetFormatter(formatter)

	log.Print("test")
	counter.Assert(t)
}

func TestSetScope(t *testing.T) {
	resetLogger()
	counter := assert.Count(1)
	formatter := NewFormatterMock(func(msg string, fields Fields) string {
		counter.Inc()
		assert.MapHas(t, fields, "version", "1.0")
		assert.MapHas(t, fields, "time", "123")
		return msg
	})
	SetFormatter(formatter)

	SetScope(Fields{"version": "1.0"})
	SetScope(Fields{"time": "123"})

	log.Print("test")
	counter.Assert(t)
}

func TestCreateWithFormatter(t *testing.T) {
	resetLogger()
	counter := assert.Count(1)
	formatter := NewFormatterMock(func(msg string, fields Fields) string {
		counter.Inc()
		return msg
	})

	logger := CreateWithFormatter(formatter)
	assert.NotNil(t, logger)

	logger.Print("test")
	counter.Assert(t)
}

func TestCreateWithScope(t *testing.T) {
	resetLogger()
	counter := assert.Count(1)
	formatter := NewFormatterMock(func(msg string, fields Fields) string {
		counter.Inc()
		assert.Equal(t, 2, len(fields))
		assert.MapHas(t, fields, "version", "1.0")
		return msg
	})
	SetFormatter(formatter)

	logger := CreateWithScope(Fields{
		"version":       "1.0",
		"formattedTime": time.Now().UTC().Format("2006-01-02"),
	})
	assert.NotNil(t, logger)

	logger.Print("test")
	counter.Assert(t)
}

func TestPrintInfo(t *testing.T) {
	printTest(t, PrintInfo)
}

func TestPrintInfof(t *testing.T) {
	printTestf(t, PrintInfof)
}

func TestPrintWarn(t *testing.T) {
	printTest(t, PrintWarn)
}

func TestPrintWarnf(t *testing.T) {
	printTestf(t, PrintWarnf)
}

func TestPrintError(t *testing.T) {
	printTest(t, PrintError)
}

func TestPrintErrorf(t *testing.T) {
	printTestf(t, PrintErrorf)
}

func TestPrintDebug(t *testing.T) {
	printTest(t, PrintDebug)
}

func TestPrintDebugf(t *testing.T) {
	printTestf(t, PrintDebugf)
}

func TestPrintTrace(t *testing.T) {
	printTest(t, PrintTrace)
}

func TestPrintTracef(t *testing.T) {
	printTestf(t, PrintTracef)
}

func TestPrintFatal(t *testing.T) {
	resetLogger()
	counter := assert.Count(1)
	formatter := NewFormatterMock(func(msg string, fields Fields) string {
		counter.Inc()
		assert.Equal(t, "test", msg)
		t.Skip()
		return msg
	})
	SetFormatter(formatter)

	PrintFatal("test")

	counter.Assert(t)
}

func TestPrintFatalf(t *testing.T) {
	resetLogger()
	counter := assert.Count(1)
	formatter := NewFormatterMock(func(msg string, fields Fields) string {
		counter.Inc()
		assert.Equal(t, "test", msg)
		t.Skip()
		return msg
	})
	SetFormatter(formatter)

	PrintFatalf("%v", "test")

	counter.Assert(t)
}

func printTest(t *testing.T, fn func(interface{})) {
	resetLogger()
	counter := assert.Count(1)
	formatter := NewFormatterMock(func(msg string, fields Fields) string {
		counter.Inc()
		assert.Equal(t, "test", msg)
		return msg
	})
	SetFormatter(formatter)

	fn("test")

	counter.Assert(t)
}

func printTestf(t *testing.T, fn func(string, ...any)) {
	resetLogger()
	counter := assert.Count(1)
	formatter := NewFormatterMock(func(msg string, fields Fields) string {
		counter.Inc()
		assert.Equal(t, "test", msg)
		return msg
	})
	SetFormatter(formatter)

	fn("%v", "test")

	counter.Assert(t)
}
