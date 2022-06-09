package logf

import (
	"goflow/tests/assert"
	"goflow/tests/mocks"
	"testing"
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

func TestMergeFields(t *testing.T) {
	tests := []struct {
		name   string
		source Fields
		fields Fields
		want   Fields
	}{
		{
			name:   "equal",
			source: Fields{"1": "1"},
			fields: Fields{"1": "1"},
			want:   Fields{"1": "1"},
		},
		{
			name:   "expanded",
			source: Fields{"1": "1"},
			fields: Fields{"2": "2"},
			want:   Fields{"1": "1", "2": "2"},
		},
		{
			name:   "change-existing-value",
			source: Fields{"1": "1"},
			fields: Fields{"1": "2", "2": "2"},
			want:   Fields{"1": "2", "2": "2"},
		},
		{
			name:   "nil-source",
			source: nil,
			fields: Fields{"1": "1"},
			want:   Fields{"1": "1"},
		},
		{
			name:   "nil-fields",
			source: Fields{"1": "1"},
			fields: nil,
			want:   Fields{"1": "1"},
		},
		{
			name:   "nil-both",
			source: nil,
			fields: nil,
			want:   Fields{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged := MergeFields(tt.source, tt.fields)
			assert.MapMatch(t, tt.want, merged)
		})
	}
}

func TestWrite(t *testing.T) {
	formatterCounter := assert.Count(1)
	writerCounter := assert.Count(1)

	formatterMock := NewFormatterMock(func(msg string, fields Fields) string {
		formatterCounter.Inc()
		return msg
	})

	writerMock := mocks.NewWriterMock(func(p []byte) (n int, err error) {
		writerCounter.Inc()
		return 0, nil
	})

	writer := formatterWriter{
		formatter: formatterMock,
		writer:    writerMock,
	}

	n, err := writer.Write([]byte("smth"))

	assert.Equal(t, 0, n)
	assert.NilError(t, err)
	formatterCounter.Assert(t)
	writerCounter.Assert(t)
}
