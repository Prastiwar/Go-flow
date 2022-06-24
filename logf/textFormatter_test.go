package logf

import (
	"testing"
)

func TestTextFormatter_Format(t *testing.T) {
	tests := []formatterTestCase{
		{
			name:     "message-without-fields",
			f:        NewTextFormatter(),
			msg:      "test",
			fields:   nil,
			expected: "test\n",
		},
		{
			name:     "message-with-fields",
			f:        NewTextFormatter(),
			msg:      "test",
			fields:   Fields{"count": 1},
			expected: "test {\"count\":1}\n",
		},
		{
			name:     "message-with-left-fields",
			f:        NewTextFormatterWith("version", "date"),
			msg:      "test",
			fields:   Fields{"count": 1, "version": "1.0", "date": "2022-12-13"},
			expected: "[1.0] [2022-12-13] test {\"count\":1}\n",
		},
	}

	testFormatter(t, tests)
}
