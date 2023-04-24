package cast_test

import (
	"strconv"
	"testing"

	"github.com/Prastiwar/Go-flow/reflection/cast"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

type customString string

type customCustomString customString

func TestAsIntToString(t *testing.T) {
	arr := []interface{}{"1", "2", "3"}

	r, ok := cast.As[string](arr)

	assert.Equal(t, true, ok)
	assert.Equal(t, len(arr), len(r))
	for i := 0; i < len(arr); i++ {
		assert.Equal(t, arr[i], r[i])
	}
}

func TestAsCustomString(t *testing.T) {
	arr := []customString{"1", "2", "3"}

	r, ok := cast.As[string](arr)

	assert.Equal(t, true, ok)
	assert.Equal(t, len(arr), len(r))
	for i := 0; i < len(arr); i++ {
		assert.Equal(t, string(arr[i]), r[i])
	}
}

func TestAsCustomCustomString(t *testing.T) {
	arr := []customCustomString{"1", "2", "3"}

	r, ok := cast.As[customString](arr)

	assert.Equal(t, true, ok)
	assert.Equal(t, len(arr), len(r))
	for i := 0; i < len(arr); i++ {
		assert.Equal(t, customString(arr[i]), r[i])
	}
}

func TestAsInvalid(t *testing.T) {
	arr := []interface{}{1, 2, 3}

	r, ok := cast.As[struct{}](arr)

	assert.Equal(t, false, ok)
	assert.Equal(t, 0, len(r))
}

func TestParseStringToInt(t *testing.T) {
	arr := []string{"1", "2", "3"}

	r, ok := cast.Parse[int32](arr)

	assert.Equal(t, true, ok)
	assert.Equal(t, len(arr), len(r))
	for i := 0; i < len(arr); i++ {
		v, err := strconv.Atoi(arr[i])
		assert.NilError(t, err)
		assert.Equal(t, int32(v), r[i])
	}
}

func TestParseIntToString(t *testing.T) {
	arr := []int32{1, 2, 3}

	r, ok := cast.Parse[string](arr)

	assert.Equal(t, true, ok)
	assert.Equal(t, len(arr), len(r))
	for i := 0; i < len(arr); i++ {
		assert.Equal(t, string(arr[i]), r[i])
	}
}

func TestParseInvalid(t *testing.T) {
	arr := []int32{1, 2, 3}

	r, ok := cast.Parse[struct{}](arr)

	assert.Equal(t, false, ok)
	assert.Equal(t, 0, len(r))
}
