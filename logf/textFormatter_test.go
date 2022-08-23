package logf

import (
	"fmt"
	"testing"
	"time"
)

func TestTextFormatter_Format(t *testing.T) {
	const monthCalendarFormat = "2006-01-02"
	now := time.Now().UTC().Format(monthCalendarFormat)

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
			name:     "message-with-left-field",
			f:        NewTextFormatterWith("version"),
			msg:      "test",
			fields:   Fields{"count": 1, "version": "1.0"},
			expected: "[1.0] test {\"count\":1}\n",
		},
		{
			name:     "message-with-time-field",
			f:        NewTextFormatter(),
			msg:      "test",
			fields:   Fields{LogTime: NewTimeField(monthCalendarFormat)},
			expected: fmt.Sprintf("[%v] test\n", now),
		},
	}

	testFormatter(t, tests)
}
