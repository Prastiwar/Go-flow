package logf

import (
	"log"
	"os"
)

const (
	Level      = "level"
	InfoLevel  = "information"
	WarnLevel  = "warning"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	DebugLevel = "debug"
	TraceLevel = "trace"
)

// resetLogger set output to os.Stderr
func resetLogger() {
	log.SetOutput(os.Stderr)
}

// SetFormatter wraps default writer with formatter and sets flags to 0
func SetFormatter(formatter Formatter) {
	out := createWrapper(log.Default(), formatter, nil)

	log.SetFlags(0)
	log.SetOutput(out)
}

// SetScope wraps default formatter with additional scope
func SetScope(fields Fields) {
	out := createWrapper(log.Default(), nil, fields)

	log.SetFlags(0)
	log.SetOutput(out)
}

// WithScope creates new instance of log.Logger based on default logger with provided fields.
// Formatter is preserved or initialized with logf.DefaultFormatter() if not set
func CreateWithScope(fields Fields) *log.Logger {
	return WithScope(log.Default(), fields)
}

// WithFormatter creates new instance of log.Logger based on default logger with replaced formatter
func CreateWithFormatter(formatter Formatter) *log.Logger {
	return WithFormatter(log.Default(), formatter)
}

func PrintError(v interface{}) {
	withLevel(log.Default(), ErrorLevel).Print(v)
}

func PrintErrorf(format string, args ...any) {
	withLevel(log.Default(), ErrorLevel).Printf(format, args...)
}

func PrintWarn(v interface{}) {
	withLevel(log.Default(), WarnLevel).Print(v)
}

func PrintWarnf(format string, args ...any) {
	withLevel(log.Default(), WarnLevel).Printf(format, args...)
}

func PrintInfo(v interface{}) {
	withLevel(log.Default(), InfoLevel).Print(v)
}

func PrintInfof(format string, args ...any) {
	withLevel(log.Default(), InfoLevel).Printf(format, args...)
}

func PrintFatal(v interface{}) {
	withLevel(log.Default(), FatalLevel).Fatal(v)
}

func PrintFatalf(format string, args ...any) {
	withLevel(log.Default(), FatalLevel).Fatalf(format, args...)
}

func PrintDebug(v interface{}) {
	withLevel(log.Default(), DebugLevel).Print(v)
}

func PrintDebugf(format string, args ...any) {
	withLevel(log.Default(), DebugLevel).Printf(format, args...)
}

func PrintTrace(v interface{}) {
	withLevel(log.Default(), TraceLevel).Print(v)
}

func PrintTracef(format string, args ...any) {
	withLevel(log.Default(), TraceLevel).Printf(format, args...)
}
