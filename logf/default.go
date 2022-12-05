package logf

import "context"

type contextLoggerKey struct{}

var loggerKey = &contextLoggerKey{}

var loggerFactory func() Logger

// From returns Logger stored in ctx using WithLogger function. In case logger was not stored
// it simply returns instance from Default() function
func From(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		return Default()
	}
	return l
}

// WithLogger returns a copy of ctx in which the Logger value is stored
func WithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// SetDefault configures behaviour of Default() function to return logger instance using specified factory
func SetDefault(factory func() Logger) {
	loggerFactory = factory
}

// Default returns logger instance by call to factory initialized using SetDefault() function.
// If it was not initialized it simply returns NewLogger()
func Default() Logger {
	if loggerFactory != nil {
		return loggerFactory()
	}
	return NewLogger()
}
