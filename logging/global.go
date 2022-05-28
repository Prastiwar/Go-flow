package logging

var defaultLog Logger = NewLogger()

func WithOptions(opts ...LogOption) {
	defaultLog = NewLogger(opts...)
}

func Info(args ...any) {
	LogInfo(defaultLog, args...)
}

func Warn(args ...any) {
	LogWarn(defaultLog, args...)
}

func Error(err error) {
	LogErr(defaultLog, err)
}
