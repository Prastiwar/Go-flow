package assert

import (
	"testing"
)

func TestMapMatch(t *testing.T) {
	tests := []struct {
		name  string
		mapA  map[string]any
		mapB  map[string]any
		fails bool
	}{
		{
			name: "success-match",
			mapA: map[string]any{
				"test":  1,
				"test2": 2,
			},
			mapB: map[string]any{
				"test2": 2,
				"test":  1,
			},
			fails: false,
		},
		{
			name: "not-match",
			mapA: map[string]any{
				"test2": 1,
				"test3": 1,
			},
			mapB: map[string]any{
				"test":  1,
				"test2": 2,
			},
			fails: true,
		},
		{
			name: "not-match-len",
			mapA: map[string]any{
				"test": 1,
			},
			mapB: map[string]any{
				"test":  1,
				"test2": 2,
			},
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			MapMatch(test, tt.mapA, tt.mapB)

			Equal(t, tt.fails, test.Failed())
		})
	}
}

func TestMapHas(t *testing.T) {
	tests := []struct {
		name  string
		m     map[string]any
		key   string
		val   any
		fails bool
	}{
		{
			name: "has",
			m: map[string]any{
				"test": 1,
			},
			key:   "test",
			val:   1,
			fails: false,
		},
		{
			name: "no-key-found",
			m: map[string]any{
				"test": 1,
			},
			key:   "test2",
			val:   1,
			fails: true,
		},
		{
			name: "value-mismatch",
			m: map[string]any{
				"test": 1,
			},
			key:   "test",
			val:   2,
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			MapHas(test, tt.m, tt.key, tt.val)

			Equal(t, tt.fails, test.Failed())
		})
	}
}
