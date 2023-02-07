// Package retry provides policy for repeating function call to handle transient errors.
package retry

import (
	"context"
	"time"
)

var _ Policy = &policy{}

// Policy is implemented by any value that has a Execute method.
// The implementation controls how to retry function and which features like
// retry count and cancel control are included.
type Policy interface {
	// Execute should call the fn at least once and try to retry it if any error occurs.
	Execute(ctx context.Context, fn func() error) error
}

// policy implements Policy interface. Controls retry execution flow and allows to
// configure retry count, wait time before retry execution and cancel predicate.
type policy struct {
	count  int
	waiter Waiter
	cancel CancelPredicate
}

// NewPolicy returns a new retry policy with configured options.
func NewPolicy(opts ...Option) *policy {
	p := &policy{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Execute calls fn at least once. It will repeat fn call until CancelPredicate will return true or
// attempts will exceed configured retry count. It will not recover or retry from panic.
func (p *policy) Execute(ctx context.Context, fn func() error) error {
	var err error
	attempts := p.count + 1
	cancel := p.cancel

	for i := 1; i <= attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}

		stop := cancel != nil && cancel(i, err)
		if stop {
			break
		}

		if err := p.wait(ctx, i, err); err != nil {
			return err
		}
	}

	return err
}

// wait retrieves wait time from configured Waiter and sleeps for given time.
// It will not wait if Waiter was not set or there will not be further attempts.
// An error is returned only when ctx returns any error.
func (p *policy) wait(ctx context.Context, attempt int, err error) error {
	if p.waiter == nil {
		return nil
	}

	if attempt == p.count+1 {
		// there is no further attempts so no need to wait
		return nil
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		return ctx.Err()
	}

	waitTime := p.waiter(attempt, err)
	waitCtx, cancel := context.WithDeadline(ctx, time.Now().Add(waitTime))
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitCtx.Done():
		return nil
	}
}
