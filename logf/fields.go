package logf

import (
	"time"
)

// Fields is used to parse scope for each log as key value pair.
type Fields map[string]interface{}

// MergeFields puts values from fields to source and returns merged Fields.
func MergeFields(source Fields, fields Fields) Fields {
	if fields == nil && source == nil {
		return make(Fields)
	}

	if fields == nil {
		return source
	}

	if source == nil {
		return fields
	}

	for k, v := range fields {
		source[k] = v
	}

	return source
}

type timeField struct {
	format string
}

// NewTimeField returns a new time field which should always return the current time
// with specified format on log formatting.
func NewTimeField(format string) *timeField {
	return &timeField{
		format: format,
	}
}

func (t *timeField) String() string {
	return time.Now().UTC().Format(t.format)
}

func (t *timeField) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}
