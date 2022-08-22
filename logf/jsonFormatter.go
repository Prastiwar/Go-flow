package logf

import "encoding/json"

// JsonFormatter implements Formatter interface. Provides message formatting as JSON output.
type JsonFormatter struct {
	pretty bool
}

// NewJsonFormatter returns a new Json Formatter with specified indentation mode.
func NewJsonFormatter(pretty bool) *JsonFormatter {
	return &JsonFormatter{pretty: pretty}
}

// Format returns string of json structure with message and fields as json object properties.
func (f *JsonFormatter) Format(msg string, fields Fields) string {
	val := MergeFields(fields, Fields{"message": msg})
	if f.pretty {
		b, _ := json.MarshalIndent(val, "", "	")
		return string(b)
	}

	b, _ := json.Marshal(val)
	return string(b)
}
