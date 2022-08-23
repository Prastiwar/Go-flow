package logf

import (
	"goflow/tests/assert"
	"goflow/tests/mocks"
	"testing"
)

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
			writerCounter := assert.Count(1)
			formatCounter := assert.Count(1)
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
		})
	}
}
