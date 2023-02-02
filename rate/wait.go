package rate

import (
	"context"
	"errors"
	"time"
)

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

// ConsumeAndWait calls Take() on Limiter and immedietely Use the Token. If ErrRateLimitExceeded error occurs
// it pauses the current goroutine until Limiter's ResetsAt() time deadline exceeds. If Take() returns any
// other error it'll immediately return this error. When ctx is canceled or ctx deadline exceeds before
// reset time it'll return this error.
func ConsumeAndWait(ctx context.Context, l Limiter) error {
	token := l.Take()
	return takeAndWait(ctx, token, token.Use())
}

// ConsumeNAndWait is extended ConsumeAndWait() function that supports BurstLimiter.TakeN().
func ConsumeNAndWait(ctx context.Context, l BurstLimiter, n uint64) error {
	token := l.TakeN(n)
	return takeAndWait(ctx, token, token.Use())
}

func takeAndWait(ctx context.Context, t Token, takeErr error) error {
	if !errors.Is(takeErr, ErrRateLimitExceeded) {
		return takeErr
	}

	if err := Wait(ctx, t.ResetsAt()); err != nil {
		return err
	}

	return nil
}
