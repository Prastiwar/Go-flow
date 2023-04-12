package memory

import (
	"context"
	"sync"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
)

var (
	_ rate.LimiterStore = (*memoryStore)(nil)
)

type memoryStore struct {
	store   sync.Map
	factory rate.LimiterAlgorithm
}

func (ls *memoryStore) Limit(ctx context.Context, key string) rate.Limiter {
	l, ok := ls.store.Load(key)
	if !ok {
		l = ls.factory()
		ls.store.Store(key, l)
	}
	return l.(rate.Limiter)
}

func (ls *memoryStore) cleanup(ctx context.Context) {
	ls.store.Range(func(key, value any) bool {
		l := value.(rate.Limiter)
		avail := l.Tokens(ctx)
		if avail >= l.Limit() {
			ls.store.Delete(key)
		}
		return true
	})
}

// NewLimiterStore returns a rate.LimiterStore which stores keys in memory using sync.Map for thread safety.
// It'll create a single goroutine to perform cleanup with provided cleanupInterval to remove unused limiters.
// Unused means when the limiter's available tokens are equal to the limit. There is no tracking for the last time used for
// a more memory-efficient solution. Adjust the cleanupInterval parameter to define how often the cleanup should be performed.
// Lowering the value means more cleanup frequency therefore more CPU usage but faster memory release.
// The cleanup time depends on cleanup execution time, meaning if the cleanup interval is set to 5s.
// It'll run cleanup on the 5th second and if cleanup execution takes 1s then the second cleanup will be performed at the 11th second.
func NewLimiterStore(ctx context.Context, alg rate.LimiterAlgorithm, cleanupInterval time.Duration) rate.LimiterStore {
	store := &memoryStore{
		factory: alg,
	}

	go func() {
		for {
			waitCtx, cancel := context.WithDeadline(ctx, time.Now().Add(cleanupInterval))
			defer cancel()

			select {
			case <-ctx.Done():
				return
			case <-waitCtx.Done():
				store.cleanup(ctx)
				continue
			}
		}
	}()

	return store
}
