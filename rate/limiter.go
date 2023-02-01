package rate

import (
	"context"
	"errors"
	"time"
)

var (
	ErrLimitReached = errors.New("Limiter rate max take count exceeded.")
)

type Token interface {
	Use() error
	Wait(ctx context.Context)
}

type Rate struct {
	MaxRate   uint64
	Remaining uint64
	ResetsAt  time.Time
}

type Limiter interface {
	Take() error

	MaxRate() uint64
	Remaining() uint64
	ResetsAt() time.Time

	Reset()
}

type BurstLimiter interface {
	Limiter

	TakeN(n uint64) error
}

// Wait pauses the current goroutine until deadline returning no error. If passed ctx is done
// before deadline - context's error is returned.
func Wait(ctx context.Context, deadline time.Time) error {
	waitCtx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitCtx.Done():
		return nil
	}
}

// TakeAndWait calls Take() on Limiter. If ErrLimitReached error occurs it pauses the current goroutine until
// Limiter's ResetsAt() time deadline. If Take() returns any other error it'll immediately return this error.
// When ctx is canceled or ctx deadline exceeds before reset time it'll return this error.
func TakeAndWait(ctx context.Context, l Limiter) error {
	return takeAndWait(ctx, l, l.Take())
}

// TakeNAndWait is extended TakeAndWait() function that supports BurstLimiter.TakeN().
func TakeNAndWait(ctx context.Context, n uint64, l BurstLimiter) error {
	return takeAndWait(ctx, l, l.TakeN(n))
}

func takeAndWait(ctx context.Context, l Limiter, takeErr error) error {
	if !errors.Is(takeErr, ErrLimitReached) {
		return takeErr
	}

	if err := Wait(ctx, l.ResetsAt()); err != nil {
		return err
	}

	return nil
}
