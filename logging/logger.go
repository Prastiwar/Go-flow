package logging

import (
	"log"
)

type Logger interface {
	Log(level Level, format string, args ...any)
	Options() LogOptions
}

type logger struct {
	options LogOptions
}

func (l *logger) Log(level Level, format string, args ...any) {
	log.Printf(format, args...)
}

func (l *logger) Options() LogOptions {
	return l.options
}

func NewLogger(opts ...LogOption) *logger {
	o := Options(opts...)

	return &logger{options: *o}
}

func LogInfo(l Logger, args ...any) {
	l.Log(Infol, format(l, l.Options().infof), args...)
}

func LogErr(l Logger, err error, args ...any) {
	reorderedArgs := []any{err}
	reorderedArgs = append(reorderedArgs, args...)
	l.Log(Errorl, format(l, l.Options().errorf), reorderedArgs...)
}

func LogWarn(l Logger, args ...any) {
	l.Log(Warnl, format(l, l.Options().warnf), args...)
}

func format(l Logger, format string) string {
	if len(format) > 0 {
		return format
	}

	return l.Options().logf
}
