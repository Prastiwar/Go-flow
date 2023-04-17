package tests

import (
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

// allocateNPointers allocates n*8 bytes of memory
func allocateNPointers(n uint64) {
	byteSlice := make([]*struct{}, 0, n)
	for i := uint64(0); i < n; i++ {
		var v struct{}
		byteSlice = append(byteSlice, &v)
	}
}

func TestMemoryUsage(t *testing.T) {
	bytes := MemoryUsage(func() { allocateNPointers(10) })

	assert.Equal(t, uint64(10*8), bytes)
}

func TestMemoryUsageFormatted(t *testing.T) {
	var bytes string

	bytes = MemoryUsageFormatted(func() { allocateNPointers(10) })
	assert.Equal(t, "80 bytes", bytes)

	bytes = MemoryUsageFormatted(func() { allocateNPointers(1024) })
	assert.Equal(t, "8 kb", bytes)
}
