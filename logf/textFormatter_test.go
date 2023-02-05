package logf_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/logf"
)

func TestTextFormatter_Format(t *testing.T) {
	const monthCalendarFormat = "2006-01-02"
	now := time.Now().UTC().Format(monthCalendarFormat)

	tests := []formatterTestCase{
		{
			name:     "message-without-fields",
			f:        logf.NewTextFormatter(),
			msg:      "test",
			fields:   nil,
			expected: "test",
		},
		{
			name:     "message-with-fields",
			f:        logf.NewTextFormatter(),
			msg:      "test",
			fields:   logf.Fields{"count": 1},
			expected: "test {\"count\":1}",
		},
		{
			name:     "message-with-left-field",
			f:        logf.NewTextFormatterWith("version"),
			msg:      "test",
			fields:   logf.Fields{"count": 1, "version": "1.0"},
			expected: "[1.0] test {\"count\":1}",
		},
		{
			name:     "message-with-time-field",
			f:        logf.NewTextFormatter(),
			msg:      "test",
			fields:   logf.Fields{logf.LogTime: logf.NewTimeField(monthCalendarFormat)},
			expected: fmt.Sprintf("[%v] test", now),
		},
	}

	testFormatter(t, tests)
}
