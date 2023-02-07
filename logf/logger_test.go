package logf_test

import (
	"testing"

	"github.com/Prastiwar/Go-flow/logf"
	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

type formatterTestCase struct {
	name     string
	f        logf.Formatter
	msg      string
	fields   logf.Fields
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
		originScope   logf.Fields
		withScope     logf.Fields
		expectedScope logf.Fields
	}{
		{
			name: "success-no-fields",
		},
		{
			name:        "success-add-field",
			originScope: logf.Fields{"time": "now"},
			withScope:   logf.Fields{"level": "1"},
			expectedScope: logf.Fields{
				"time":  "now",
				"level": "1",
			},
		},
		{
			name:          "success-override-field",
			originScope:   logf.Fields{"time": "now"},
			withScope:     logf.Fields{"time": "today"},
			expectedScope: logf.Fields{"time": "today"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerMock := logf.NewLogger(logf.WithFields(tt.originScope))
			logger := logf.WithScope(loggerMock, tt.withScope)
			assert.MapMatch(t, tt.expectedScope, logger.Scope())

			assert.NotNil(t, logger)
		})
	}
}

func TestPrinting(t *testing.T) {
	tests := []struct {
		name  string
		print func(logf.Logger, string, ...any)
	}{
		{
			name: "success-info",
			print: func(l logf.Logger, s string, a ...any) {
				l.Info(s)
			},
		},
		{
			name: "success-infof",
			print: func(l logf.Logger, s string, a ...any) {
				l.Infof(s, a...)
			},
		},
		{
			name: "success-error",
			print: func(l logf.Logger, s string, a ...any) {
				l.Error(s)
			},
		},
		{
			name: "success-errorf",
			print: func(l logf.Logger, s string, a ...any) {
				l.Errorf(s, a...)
			},
		},
		{
			name: "success-debug",
			print: func(l logf.Logger, s string, a ...any) {
				l.Debug(s)
			},
		},
		{
			name: "success-debugf",
			print: func(l logf.Logger, s string, a ...any) {
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
			formatMock := mocks.FormatterMock{
				OnFormat: func(msg string, fields logf.Fields) string {
					formatCounter.Inc()
					return msg
				},
			}
			loggerMock := logf.NewLogger(
				logf.WithOutput(writerMock),
				logf.WithFormatter(formatMock),
			)

			tt.print(loggerMock, "%v", "test")

			writerCounter.Assert(t)
			formatCounter.Assert(t)
		})
	}
}
