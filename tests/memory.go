package tests

import (
	"fmt"
	"math"
	"runtime"
)

// MemoryUsage returns memory usage in bytes.
func MemoryUsage(f func()) uint64 {
	var m1, m2 runtime.MemStats
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	runtime.ReadMemStats(&m1)

	f()

	runtime.ReadMemStats(&m2)
	return m2.TotalAlloc - m1.TotalAlloc
}

// MemoryUsageFormatted returns memory usage in highest possible unit as string.
// 10015 bytes will be returned as "10 kb"
func MemoryUsageFormatted(f func()) string {
	mem := MemoryUsage(f)
	unitSize, unitName := getUnit(mem)
	return fmt.Sprintf("%v %v", unitSize, unitName)
}

func getUnit(mem uint64) (uint64, string) {
	units := []string{
		"tb", "gb", "mb", "kb",
	}

	boundaryIndex := float64(len(units))
	for _, v := range units {
		boundary := uint64(math.Pow(1024, boundaryIndex))
		if mem > boundary {
			return mem / boundary, v
		}
		boundaryIndex--
	}

	return mem, "bytes"
}
