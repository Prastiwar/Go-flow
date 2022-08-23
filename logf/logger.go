package logf

import (
	"fmt"
	"io"
	"strings"
)

const (
	// LogTime is shared key used for time field
	LogTime = "log_time"

	// Level is shared key used for level field
	Level = "level"

	InfoLevel  = "INFO"
	ErrorLevel = "ERR"
	DebugLevel = "DEBUG"
)

type Formatter interface {
	Format(msg string, fields Fields) string
}

type Logger interface {
	Output() io.Writer
	Scope() Fields
	Formatter() Formatter

	Error(v interface{})
	Errorf(format string, args ...any)

	Info(v interface{})
	Infof(format string, args ...any)

	Debug(v interface{})
	Debugf(format string, args ...any)
}

type wrappedLogger struct {
	writer    io.Writer
	formatter Formatter
	fields    Fields
}

// NewLogger returns a new Logger which is wrapper for new log.Logger.
func NewLogger(opts ...LoggerOption) Logger {
	options := NewLoggerOptions(opts...)
	return &wrappedLogger{
		writer:    options.Output,
		formatter: options.Format,
		fields:    options.Scope,
	}
}

// WithScope creates new instance of log.Logger with provided fields.
// Formatter is preserved or initialized with logf.DefaultFormatter() if not set
func WithScope(logger Logger, fields Fields) Logger {
	return &wrappedLogger{
		writer:    logger.Output(),
		formatter: logger.Formatter(),
		fields:    MergeFields(logger.Scope(), fields),
	}
}

func (l *wrappedLogger) Output() io.Writer {
	return l.writer
}

func (l *wrappedLogger) Formatter() Formatter {
	return l.formatter
}

func (l *wrappedLogger) Scope() Fields {
	return l.fields
}

func (l *wrappedLogger) Error(v interface{}) {
	l.print(ErrorLevel, v)
}

func (l *wrappedLogger) Errorf(format string, args ...any) {
	l.printf(ErrorLevel, format, args...)
}

func (l *wrappedLogger) Info(v interface{}) {
	l.print(InfoLevel, v)
}

func (l *wrappedLogger) Infof(format string, args ...any) {
	l.printf(InfoLevel, format, args...)
}

func (l *wrappedLogger) Debug(v interface{}) {
	l.print(DebugLevel, v)
}

func (l *wrappedLogger) Debugf(format string, args ...any) {
	l.printf(DebugLevel, format, args...)
}

func (l *wrappedLogger) print(level string, v interface{}) {
	formattedMsg := l.formatMessage(level, fmt.Sprint(v))
	l.output(formattedMsg)
}

func (l *wrappedLogger) printf(level string, format string, args ...any) {
	formattedMsg := l.formatMessage(level, fmt.Sprintf(format, args...))
	l.output(formattedMsg)
}

func (l *wrappedLogger) formatMessage(level string, message string) string {
	levelField := Fields{Level: level}
	fields := MergeFields(l.fields, levelField)
	msg := strings.TrimSuffix(message, "\n")
	return l.formatter.Format(msg, fields)
}

func (l *wrappedLogger) output(message string) {
	buf := []byte(message + "\n")
	_, _ = l.writer.Write(buf)
}
