package logf

type JsonFormatter struct {
	pretty bool
}

func (f *JsonFormatter) Format(msg string, fields Fields) string {
	// TODO: implement
	if f.pretty {

	} else {

	}

	return msg
}
