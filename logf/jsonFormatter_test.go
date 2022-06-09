package logf

import (
	"testing"
)

func TestJsonFormatter_Format(t *testing.T) {
	prettyFormatter := NewJsonFormatter(true)
	compactFormatter := NewJsonFormatter(false)
	fields := Fields{"count": 1, "version": "1.0"}

	tests := []formatterTestCase{
		{
			name:     "compact-message-without-fields",
			f:        compactFormatter,
			msg:      "test",
			fields:   nil,
			expected: "{\"message\":\"test\"}",
		},
		{
			name:   "pretty-message-without-fields",
			f:      prettyFormatter,
			msg:    "test",
			fields: nil,
			expected: `{
	"message": "test"
}`,
		},
		{
			name:     "compact-message-with-fields",
			f:        compactFormatter,
			msg:      "test",
			fields:   fields,
			expected: "{\"count\":1,\"message\":\"test\",\"version\":\"1.0\"}",
		},
		{
			name:   "pretty-message-with-fields",
			f:      prettyFormatter,
			msg:    "test",
			fields: fields,
			expected: `{
	"count": 1,
	"message": "test",
	"version": "1.0"
}`,
		},
	}

	testFormatter(t, tests)
}
