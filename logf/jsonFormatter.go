package logf

import "encoding/json"

type JsonFormatter struct {
	pretty bool
}

func NewJsonFormatter(pretty bool) *JsonFormatter {
	return &JsonFormatter{pretty: pretty}
}

func (f *JsonFormatter) Format(msg string, fields Fields) string {
	val := MergeFields(fields, Fields{"message": msg})
	if f.pretty {
		b, _ := json.MarshalIndent(val, "", "	")
		return string(b)
	}

	b, _ := json.Marshal(val)
	return string(b)
}
