package logf_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/logf"
)

func TestJsonFormatter_Format(t *testing.T) {
	const monthCalendarFormat = "2006-01-02"
	now := time.Now().UTC().Format(monthCalendarFormat)
	prettyFormatter := logf.NewJsonFormatter(true)
	compactFormatter := logf.NewJsonFormatter(false)
	fields := logf.Fields{"count": 1, "version": "1.0"}

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
		{
			name:     "compact-message-with-time-field",
			f:        compactFormatter,
			msg:      "test",
			fields:   logf.Fields{logf.LogTime: logf.NewTimeField(monthCalendarFormat)},
			expected: fmt.Sprintf("{\"%v\":\"%v\",\"message\":\"test\"}", logf.LogTime, now),
		},
	}

	testFormatter(t, tests)
}
