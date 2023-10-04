package assert_test

import (
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestElementsMatch(t *testing.T) {
	tests := []struct {
		name  string
		arrA  []any
		arrB  []any
		fails bool
	}{
		{
			name: "matches",
			arrA: []any{
				1, 2, 3,
			},
			arrB: []any{
				3, 2, 1,
			},
			fails: false,
		},
		{
			name: "not-matches-similar",
			arrA: []any{
				1, 2, 3,
			},
			arrB: []any{
				1, 2, 3, 1,
			},
			fails: true,
		},
		{
			name: "not-matches-totally-different",
			arrA: []any{
				1, 2, 3, 6, 7,
			},
			arrB: []any{
				3, 2, 1, 4, 6, 6,
			},
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			assert.ElementsMatch(test, tt.arrA, tt.arrB)

			assert.Equal(t, tt.fails, test.Failed())
		})
	}
}
