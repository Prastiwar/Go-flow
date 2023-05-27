package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
)

var (
	_ rate.LimiterStore = (*memoryStore)(nil)
)

var (
	ErrCleanupFailure   = errors.New("at least one error occured at cleanup")
	ErrMissingAlgorithm = errors.New("nil LimiterAlgorithm was passed to constructor")
)

type memoryStore struct {
	store   sync.Map
	factory rate.LimiterAlgorithm
}

func (ls *memoryStore) Limit(ctx context.Context, key string) (rate.Limiter, error) {
	l, ok := ls.store.Load(key)
	if !ok {
		l = ls.factory()
		ls.store.Store(key, l)
	}
	return l.(rate.Limiter), nil
}

func (ls *memoryStore) cleanup(ctx context.Context) error {
	var errs error

	ls.store.Range(func(key, value any) bool {
		l := value.(rate.Limiter)
		avail, err := l.Tokens(ctx)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("unexpected error for Limiter with key '%v': %w", key, err))
			return true
		}
		if avail >= l.Limit() {
			ls.store.Delete(key)
		}
		return true
	})

	if errs != nil {
		errs = errors.Join(ErrCleanupFailure, errs)
	}
	return errs
}

type Options struct {
	ErrorHandler func(err error)
}

type Option func(o *Options)

func WithErrorHandler(errorHandler func(err error)) Option {
	return func(o *Options) {
		o.ErrorHandler = errorHandler
	}
}

func NewOptions(opts ...Option) *Options {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// NewLimiterStore returns a rate.LimiterStore which stores keys in memory using sync.Map for thread safety.
// It'll create a single goroutine to perform cleanup with provided cleanupInterval to remove unused limiters.
// If cleanupInterval is less or equal 0, the cleanup goroutine will not run and no cleanup will ever be available for this instance.
// Unused means when the limiter's available tokens are equal to the limit. There is no tracking for the last time used for
// a more memory-efficient solution. Adjust the cleanupInterval parameter to define how often the cleanup should be performed.
// Lowering the value means more cleanup frequency therefore more CPU usage but faster memory release.
// The cleanup time depends on cleanup execution time, meaning if the cleanup interval is set to 5s.
// It'll run cleanup on the 5th second and if cleanup execution takes 1s then the second cleanup will be performed at the 11th second.
func NewLimiterStore(ctx context.Context, alg rate.LimiterAlgorithm, cleanupInterval time.Duration, opts ...Option) (rate.LimiterStore, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if alg == nil {
		return nil, ErrMissingAlgorithm
	}

	options := NewOptions(opts...)

	store := &memoryStore{
		factory: alg,
	}

	if cleanupInterval > time.Duration(0) {
		go func() {
			for {
				waitCtx, cancel := context.WithDeadline(ctx, time.Now().Add(cleanupInterval))
				defer cancel()

				select {
				case <-ctx.Done():
					return
				case <-waitCtx.Done():
					err := store.cleanup(ctx)
					if err != nil && options.ErrorHandler != nil {
						options.ErrorHandler(err)
					}
					continue
				}
			}
		}()
	}

	return store, nil
}
