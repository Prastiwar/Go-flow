package logf

import (
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

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
