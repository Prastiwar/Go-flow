package cast

import (
	"goflow/tests/assert"
	"testing"
)

type customString string

func TestAsString(t *testing.T) {
	arr := []interface{}{"1", "2", "3"}

	r, ok := As[string](arr)

	assert.Equal(t, true, ok)
	assert.Equal(t, len(arr), len(r))
	for i := 0; i < len(arr); i++ {
		assert.Equal(t, arr[i], r[i])
	}
}

func TestAsInvalid(t *testing.T) {
	arr := []customString{"1", "2", "3"}

	r, ok := As[string](arr)

	assert.Equal(t, false, ok)
	assert.Equal(t, 0, len(r))
}