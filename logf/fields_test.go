package logf_test

import (
	"testing"

	"github.com/Prastiwar/Go-flow/logf"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestMergeFields(t *testing.T) {
	tests := []struct {
		name   string
		source logf.Fields
		fields logf.Fields
		want   logf.Fields
	}{
		{
			name:   "equal",
			source: logf.Fields{"1": "1"},
			fields: logf.Fields{"1": "1"},
			want:   logf.Fields{"1": "1"},
		},
		{
			name:   "expanded",
			source: logf.Fields{"1": "1"},
			fields: logf.Fields{"2": "2"},
			want:   logf.Fields{"1": "1", "2": "2"},
		},
		{
			name:   "change-existing-value",
			source: logf.Fields{"1": "1"},
			fields: logf.Fields{"1": "2", "2": "2"},
			want:   logf.Fields{"1": "2", "2": "2"},
		},
		{
			name:   "nil-source",
			source: nil,
			fields: logf.Fields{"1": "1"},
			want:   logf.Fields{"1": "1"},
		},
		{
			name:   "nil-fields",
			source: logf.Fields{"1": "1"},
			fields: nil,
			want:   logf.Fields{"1": "1"},
		},
		{
			name:   "nil-both",
			source: nil,
			fields: nil,
			want:   logf.Fields{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged := logf.MergeFields(tt.source, tt.fields)
			assert.MapMatch(t, tt.want, merged)
		})
	}
}
