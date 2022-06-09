package logf

import "log"

func createWrapper(logger *log.Logger, formatter Formatter, fields Fields) formatterWriter {
	writer := logger.Writer()
	writerf, ok := writer.(formatterWriter)
	if ok {
		writer = writerf.writer
		if formatter == nil {
			formatter = writerf.formatter
		}
		fields = MergeFields(writerf.fields, fields)
	}

	if formatter == nil {
		formatter = DefaultFormatter()
	}

	out := formatterWriter{
		formatter: formatter,
		writer:    writer,
		fields:    fields,
	}

	return out
}

// WithScope creates new instance of log.Logger with provided fields.
// Formatter is preserved or initialized with logf.DefaultFormatter() if not set
func WithScope(logger *log.Logger, fields Fields) *log.Logger {
	out := createWrapper(logger, nil, fields)
	l := log.New(out, "", 0)
	return l
}

// WithFormatter creates new instance of log.Logger based on parent logger with replaced formatter
func WithFormatter(logger *log.Logger, formatter Formatter) *log.Logger {
	out := createWrapper(logger, formatter, nil)
	l := log.New(out, "", 0)
	return l
}

func Error(logger *log.Logger, v interface{}) {
	withLevel(logger, ErrorLevel).Print(v)
}

func Errorf(logger *log.Logger, format string, args ...any) {
	withLevel(logger, ErrorLevel).Printf(format, args...)
}

func Warn(logger *log.Logger, v interface{}) {
	withLevel(logger, WarnLevel).Print(v)
}

func Warnf(logger *log.Logger, format string, args ...any) {
	withLevel(logger, WarnLevel).Printf(format, args...)
}

func Info(logger *log.Logger, v interface{}) {
	withLevel(logger, InfoLevel).Print(v)
}

func Infof(logger *log.Logger, format string, args ...any) {
	withLevel(logger, InfoLevel).Printf(format, args...)
}

func Fatal(logger *log.Logger, v interface{}) {
	withLevel(logger, FatalLevel).Fatal(v)
}

func Fatalf(logger *log.Logger, format string, args ...any) {
	withLevel(logger, FatalLevel).Fatalf(format, args...)
}

func Debug(logger *log.Logger, v interface{}) {
	withLevel(logger, DebugLevel).Print(v)
}

func Debugf(logger *log.Logger, format string, args ...any) {
	withLevel(logger, DebugLevel).Printf(format, args...)
}

func Trace(logger *log.Logger, v interface{}) {
	withLevel(logger, TraceLevel).Print(v)
}

func Tracef(logger *log.Logger, format string, args ...any) {
	withLevel(logger, TraceLevel).Printf(format, args...)
}

func withLevel(logger *log.Logger, level string) *log.Logger {
	return WithScope(logger, Fields{Level: level})
}
