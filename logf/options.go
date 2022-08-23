package logf

import (
	"io"
	"log"
)

// LoggerOptions is set of optional parameters which can configure new Logger.
type LoggerOptions struct {
	Output io.Writer
	Format Formatter
	Scope  Fields
}

// LoggerOption is function which mutates the LoggerOptions specific field.
type LoggerOption func(*LoggerOptions)

// NewLoggerOptions returns a new LoggerOptions.
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

// WithOutput returns a new LoggerOption to configure log writer output.
func WithOutput(w io.Writer) LoggerOption {
	return func(o *LoggerOptions) {
		o.Output = w
	}
}

// WithFormatter returns a new LoggerOption to configure log message formatter.
func WithFormatter(f Formatter) LoggerOption {
	return func(o *LoggerOptions) {
		o.Format = f
	}
}

// WithFormatter returns a new LoggerOption to configure log initial scope.
func WithFields(f Fields) LoggerOption {
	return func(o *LoggerOptions) {
		o.Scope = f
	}
}
