package assert

import "testing"

func ElementsMatch[T any](t *testing.T, arrA, arrB []T) {
	if len(arrA) != len(arrB) {
		t.Errorf("expected same slice length: arrA: '%v', arrB: '%v'", len(arrA), len(arrB))
	}

	// TODO: assert for elemenets
}
