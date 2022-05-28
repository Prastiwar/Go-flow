package logging

type LogMock struct {
	options LogOptions
	logFunc func(level Level, format string, args ...any)
}

func NewLogMock(logFunc func(level Level, format string, args ...any), opts ...LogOption) *LogMock {
	return &LogMock{
		options: *Options(opts...),
		logFunc: logFunc,
	}
}

func (l *LogMock) Log(level Level, format string, args ...any) {
	l.logFunc(level, format, args...)
}

func (l *LogMock) Options() LogOptions {
	return l.options
}
