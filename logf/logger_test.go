package logf

import (
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

type formatterTestCase struct {
	name     string
	f        Formatter
	msg      string
	fields   Fields
	expected string
}

func testFormatter(t *testing.T, tests []formatterTestCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.f.Format(tt.msg, tt.fields)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWithScope(t *testing.T) {
	tests := []struct {
		name          string
		originScope   Fields
		withScope     Fields
		expectedScope Fields
	}{
		{
			name: "success-no-fields",
		},
		{
			name:        "success-add-field",
			originScope: Fields{"time": "now"},
			withScope:   Fields{"level": "1"},
			expectedScope: Fields{
				"time":  "now",
				"level": "1",
			},
		},
		{
			name:          "success-override-field",
			originScope:   Fields{"time": "now"},
			withScope:     Fields{"time": "today"},
			expectedScope: Fields{"time": "today"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerMock := NewLogger(WithFields(tt.originScope))
			logger := WithScope(loggerMock, tt.withScope)
			assert.MapMatch(t, tt.expectedScope, logger.Scope())

			assert.NotNil(t, logger)
		})
	}
}

func TestPrinting(t *testing.T) {
	tests := []struct {
		name  string
		print func(Logger, string, ...any)
	}{
		{
			name: "success-info",
			print: func(l Logger, s string, a ...any) {
				l.Info(s)
			},
		},
		{
			name: "success-infof",
			print: func(l Logger, s string, a ...any) {
				l.Infof(s, a...)
			},
		},
		{
			name: "success-error",
			print: func(l Logger, s string, a ...any) {
				l.Error(s)
			},
		},
		{
			name: "success-errorf",
			print: func(l Logger, s string, a ...any) {
				l.Errorf(s, a...)
			},
		},
		{
			name: "success-debug",
			print: func(l Logger, s string, a ...any) {
				l.Debug(s)
			},
		},
		{
			name: "success-debugf",
			print: func(l Logger, s string, a ...any) {
				l.Debugf(s, a...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writerCounter := assert.Count(t, 1)
			formatCounter := assert.Count(t, 1)
			writerMock := mocks.NewWriterMock(func(p []byte) (n int, err error) {
				writerCounter.Inc()
				return 0, nil
			})
			formatMock := NewFormatterMock(func(msg string, fields Fields) string {
				formatCounter.Inc()
				return msg
			})
			loggerMock := NewLogger(
				WithOutput(writerMock),
				WithFormatter(formatMock),
			)

			tt.print(loggerMock, "%v", "test")

			writerCounter.Assert(t)
			formatCounter.Assert(t)
		})
	}
}
