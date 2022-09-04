package assert

import "testing"

// MapMatch asserts all keys and corresponding values are applied to both map.
func MapMatch[K comparable, V any](t *testing.T, mapA, mapB map[K]V) {
	if len(mapA) != len(mapB) {
		t.Errorf("expected same map length: mapA: '%v', mapB: '%v'", len(mapA), len(mapB))
	}

	for k, v := range mapA {
		bV, ok := mapB[k]
		if !ok {
			t.Errorf("not found key: '%v' in mapB", k)
		} else {
			if any(v) != any(bV) {
				t.Errorf("key: '%v', mapA value: '%v', mapB value: '%v'", k, v, bV)
			}
		}
	}
}

func MapHas[K comparable, V any](t *testing.T, m map[K]V, key K, val V) {
	v, ok := m[key]
	if !ok {
		t.Errorf("not found key: '%v'", key)
	}

	if any(v) != any(val) {
		Equal(t, val, v)
	}
}
