package logging

type Level int

const (
	Infol Level = iota
	Errorl
	Warnl
)

type LogOptions struct {
	logf   string
	infof  string
	warnf  string
	errorf string
}

type LogOption func(*LogOptions)

func DefaultOptions() *LogOptions {
	return &LogOptions{
		logf: "",
	}
}

func Options(opts ...LogOption) *LogOptions {
	o := DefaultOptions()

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithLogFormat(format string) LogOption {
	return func(o *LogOptions) {
		o.logf = format
	}
}

func WithInfoFormat(format string) LogOption {
	return func(o *LogOptions) {
		o.infof = format
	}
}

func WithWarnFormat(format string) LogOption {
	return func(o *LogOptions) {
		o.warnf = format
	}
}

func WithErrorFormat(format string) LogOption {
	return func(o *LogOptions) {
		o.errorf = format
	}
}
