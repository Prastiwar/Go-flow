package logf

import (
	"context"
)

type contextKeyValue struct{}

var ctxKey = &contextKeyValue{}

// WithLogger returns a copy of parent in which the logger value is stored.
// This is useful if want to pass scoped logger across API request within context
func WithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctxKey, l)
}

// Logger gets logger instance stored within given context. It'll panic if context does not contian logger
func Logger(ctx context.Context) Logger {
	return ctx.Value(ctxKey).(Logger)
}