package memory_test

import (
	"context"
	"errors"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/rate/memory"
	"github.com/Prastiwar/Go-flow/rate/slidingwindow"
	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

// newStore is a shorthand to return rate.NewMemoryLimiterStore with slidingwindow.SlidingWindowCounter
func newStore(ctx context.Context, cleanupInterval time.Duration, opts ...memory.Option) rate.LimiterStore {
	slidingWindow, err := slidingwindow.NewAlgorithm(40, time.Minute, 20)
	if err != nil {
		panic(err)
	}

	store, err := memory.NewLimiterStore(ctx, slidingWindow, cleanupInterval, opts...)
	if err != nil {
		panic(err)
	}
	return store
}

func TestGoroutineCount(t *testing.T) {
	// Arrange
	goroutinesCount := runtime.NumGoroutine()
	ctx, cancel := context.WithCancel(context.Background())

	// Act
	store := newStore(ctx, time.Second)
	for i := 0; i < 10; i++ {
		_, _ = store.Limit(context.Background(), strconv.Itoa(i))
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
	ctx := context.Background()

	store := newStore(ctx, time.Second/2)

	// Act - store
	expected, err := store.Limit(ctx, key)

	// Assert - store
	assert.NilError(t, err)

	actual, err := store.Limit(ctx, key)
	assert.NilError(t, err)
	assert.Equal(t, expected, actual, "equal before cleanup")

	// Act - cleanup
	time.Sleep(time.Second)

	// Assert - cleanup
	actual, err = store.Limit(ctx, key)
	assert.NilError(t, err)
	assert.NotEqual(t, expected, actual, "different after cleanup")
}

func TestErrorHandling(t *testing.T) {
	// Arrange
	const key = "{userId}"
	ctx := context.Background()

	onTokensCaller := assert.Count(t, 1)
	onErrorCaller := assert.Count(t, 1)
	limiterMock := mocks.LimiterMock{
		OnTokens: func(ctx context.Context) (uint64, error) {
			onTokensCaller.Inc()
			return 0, errors.New("invalid-on-tokens")
		},
	}
	onError := func(err error) {
		onErrorCaller.Inc()
		assert.ErrorWith(t, err, "invalid-on-tokens")
	}

	store, err := memory.NewLimiterStore(ctx, func() rate.Limiter { return limiterMock }, time.Second/3, memory.WithErrorHandler(onError))
	assert.NilError(t, err)

	_, err = store.Limit(ctx, key)
	assert.NilError(t, err)

	// Act
	time.Sleep(time.Second / 2)

	// Assert
	onTokensCaller.Assert(t, "expected Limiter.Tokens to be called once")
	onErrorCaller.Assert(t, "expected error handler to be called once")
}

func TestNewLimiterStoreValidation(t *testing.T) {
	t.Run("nil-algorithm", func(t *testing.T) {
		_, err := memory.NewLimiterStore(nil, nil, time.Duration(0))
		assert.ErrorIs(t, err, memory.ErrMissingAlgorithm)
	})

	t.Run("zero-goroutine", func(t *testing.T) {
		goroutinesCount := runtime.NumGoroutine()
		_ = newStore(nil, time.Duration(0))
		assert.Equal(t, goroutinesCount, runtime.NumGoroutine())
	})
}
