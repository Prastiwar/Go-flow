package tests

import (
	"fmt"
	"runtime"
)

func MemoryUsage(f func()) uint64 {
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	f()

	runtime.ReadMemStats(&m2)
	return m2.TotalAlloc - m1.TotalAlloc
}

func MemoryUsageFormat(f func()) string {
	mem := MemoryUsage(f)
	sizeName := "bytes"

	switch {
	case mem > 1024:
		sizeName = "kb"
		mem = mem / 1024
	case mem > 1024*1024:
		sizeName = "mb"
		mem = mem / (1024 * 1024)
	case mem > 1024*1024*1024:
		sizeName = "gb"
		mem = mem / (1024 * 1024 * 1024)
	case mem > 1024*1024*1024*1024:
		sizeName = "tb"
		mem = mem / (1024 * 1024 * 1024 * 1024)
	}

	return fmt.Sprintf("%v %v", mem, sizeName)
}
