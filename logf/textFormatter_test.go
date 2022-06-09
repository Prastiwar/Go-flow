package logf

import (
	"testing"
)

func TestTextFormatter_Format(t *testing.T) {
	formatter := NewTextFormatter()
	fields := Fields{"count": 1, "version": "1.0"}

	tests := []formatterTestCase{
		{
			name:     "compact-message-without-fields",
			f:        formatter,
			msg:      "test",
			fields:   nil,
			expected: "test",
		},
		{
			name:     "pretty-message-without-fields",
			f:        formatter,
			msg:      "test",
			fields:   nil,
			expected: "test",
		},
		{
			name:     "compact-message-with-fields",
			f:        formatter,
			msg:      "test",
			fields:   fields,
			expected: "test",
		},
		{
			name:     "pretty-message-with-fields",
			f:        formatter,
			msg:      "test",
			fields:   fields,
			expected: "test",
		},
	}

	testFormatter(t, tests)
}
