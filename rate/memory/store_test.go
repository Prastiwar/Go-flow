package memory_test

import (
	"context"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/rate/memory"
	"github.com/Prastiwar/Go-flow/rate/slidingwindow"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

// newStore is a shorthand to return rate.NewMemoryLimiterStore with slidingwindow.SlidingWindowCounter
func newStore(ctx context.Context, cleanupInterval time.Duration) rate.LimiterStore {
	slidingWindow, err := slidingwindow.NewAlgorithm(40, time.Minute, 20)
	if err != nil {
		panic(err)
	}

	return memory.NewLimiterStore(ctx, slidingWindow, cleanupInterval)
}

func TestGoroutineCount(t *testing.T) {
	// Arrange
	goroutinesCount := runtime.NumGoroutine()
	ctx, cancel := context.WithCancel(context.Background())

	// Act
	store := newStore(ctx, time.Second)
	for i := 0; i < 10; i++ {
		_ = store.Limit(strconv.Itoa(i))
	}

	// Assert
	assert.Equal(t, goroutinesCount+1, runtime.NumGoroutine())

	cancel()
	runtime.GC()
	assert.Equal(t, goroutinesCount, runtime.NumGoroutine())
}

func TestCleanup(t *testing.T) {
	// Arrange
	const key = "{userId}"
	store := newStore(context.Background(), time.Second/2)

	// Act
	expected := store.Limit(key)

	// Assert
	actual := store.Limit(key)
	assert.Equal(t, expected, actual, "equal before cleanup")

	time.Sleep(time.Second)

	actual = store.Limit(key)
	assert.NotEqual(t, expected, actual, "different after cleanup")
}
