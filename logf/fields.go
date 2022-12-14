package logf

import (
	"time"
)

// Fields is used to parse scope for each log as key value pair.
type Fields map[string]interface{}

// MergeFields returns a new instace of Fields which contains values from source and upserts fields to it.
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

	result := make(Fields, len(source)+len(fields))

	// copy source fields
	for k, v := range source {
		result[k] = v
	}

	// upsert new fields
	for k, v := range fields {
		result[k] = v
	}

	return result
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
