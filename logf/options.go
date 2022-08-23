package logf

import (
	"io"
	"log"
)

type LoggerOptions struct {
	Output io.Writer
	Format Formatter
	Scope  Fields
}

type LoggerOption func(*LoggerOptions)

func NewLoggerOptions(opts ...LoggerOption) *LoggerOptions {
	o := &LoggerOptions{}
	for _, opt := range opts {
		opt(o)
	}

	if o.Output == nil {
		o.Output = log.Default().Writer()
	}

	if o.Format == nil {
		o.Format = NewTextFormatter()
	}

	return o
}

func WithOutput(w io.Writer) LoggerOption {
	return func(o *LoggerOptions) {
		o.Output = w
	}
}

func WithFormatter(f Formatter) LoggerOption {
	return func(o *LoggerOptions) {
		o.Format = f
	}
}

func WithFields(f Fields) LoggerOption {
	return func(o *LoggerOptions) {
		o.Scope = f
	}
}
