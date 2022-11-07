package assert

import (
	"fmt"
	"testing"
)

// MapMatch asserts all keys and corresponding values are applied to both map.
func MapMatch[K comparable, V any](t *testing.T, mapA, mapB map[K]V, prefixes ...string) {
	t.Helper()
	if len(mapA) != len(mapB) {
		errorf(t, fmt.Sprintf("expected same map length: mapA: '%v', mapB: '%v'", len(mapA), len(mapB)), prefixes...)
	}

	for k, v := range mapA {
		bV, ok := mapB[k]
		if !ok {
			errorf(t, fmt.Sprintf("not found key: '%v' in mapB", k), prefixes...)
		} else {
			if any(v) != any(bV) {
				errorf(t, fmt.Sprintf("key: '%v', mapA value: '%v', mapB value: '%v'", k, v, bV), prefixes...)
			}
		}
	}
}

// MapHas asserts that map contains specified key with equal value.
func MapHas[K comparable, V any](t *testing.T, m map[K]V, key K, val V, prefixes ...string) {
	t.Helper()
	v, ok := m[key]
	if !ok {
		errorf(t, fmt.Sprintf("not found key: '%v'", key), prefixes...)
	}

	if any(v) != any(val) {
		Equal(t, val, v)
	}
}
