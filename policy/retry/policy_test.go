package retry_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/policy/retry"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestPolicy_Execute(t *testing.T) {
	const (
		firstAttemptTime  = 100 * time.Millisecond
		secondAttemptTime = 200 * time.Millisecond
		deltaDivider      = 2
	)

	var timeStart time.Time

	withTwoAttemptsWaitAssertion := retry.WithCancelPredicate(func(attempt int, err error) bool {
		if attempt == 1 {
			// can't check it now - start timer
			timeStart = time.Now()
			return false
		}

		waited := time.Since(timeStart)
		// check first attempt wait time
		if attempt == 2 {
			assert.Approximately(t, firstAttemptTime, waited, firstAttemptTime/deltaDivider, "first attempt wait time failed")
		} else if attempt == 3 {
			// check second attempt wait time
			assert.Approximately(t, secondAttemptTime, waited, secondAttemptTime/deltaDivider, "second attempt wait time failed")
		}
		timeStart = time.Now()
		return false
	})

	tests := []struct {
		name     string
		ctx      func(t *testing.T) context.Context
		p        func(t *testing.T) retry.Policy
		fn       func(attempt int) error
		asserter assert.ResultErrorFunc[int]
	}{
		{
			name: "success-after-single-retry",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			p: func(t *testing.T) retry.Policy {
				c := assert.Count(t, 1, "failed retry call")
				return retry.NewPolicy(
					retry.WithCount(1),
					retry.WithCancelPredicate(func(attempt int, err error) bool {
						c.Inc()
						return false
					}),
				)
			},
			fn: func(attempt int) error {
				if attempt == 1 {
					return errors.New("invalid")
				}
				return nil
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.NilError(t, err)
				assert.Equal(t, attempt, 2)
			},
		},
		{
			name: "success-3-retries-cancel",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			p: func(t *testing.T) retry.Policy {
				c := assert.Count(t, 3, "failed retry call")
				return retry.NewPolicy(
					retry.WithCount(5),
					retry.WithCancelPredicate(func(attempt int, err error) bool {
						c.Inc()
						return attempt == 3
					}),
				)
			},
			fn: func(attempt int) error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.Error(t, err)
				assert.Equal(t, attempt, 3)
			},
		},
		{
			name: "success-no-cancel-no-error",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			p: func(t *testing.T) retry.Policy {
				return retry.NewPolicy(
					retry.WithCount(5),
				)
			},
			fn: func(attempt int) error {
				if attempt == 1 {
					return errors.New("invalid")
				}
				return nil
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.NilError(t, err)
				assert.Equal(t, attempt, 2)
			},
		},
		{
			name: "success-no-cancel-no-error-no-context",
			ctx: func(t *testing.T) context.Context {
				return nil
			},
			p: func(t *testing.T) retry.Policy {
				return retry.NewPolicy(
					retry.WithCount(5),
					retry.WithWaitTimes(firstAttemptTime, secondAttemptTime),
				)
			},
			fn: func(attempt int) error {
				if attempt == 1 {
					return errors.New("invalid")
				}
				return nil
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.NilError(t, err)
				assert.Equal(t, attempt, 2)
			},
		},
		{
			name: "success-context-canceled",
			ctx: func(t *testing.T) context.Context {
				return canceledContext()
			},
			p: func(t *testing.T) retry.Policy {
				return retry.NewPolicy(
					retry.WithCount(3),
					retry.WithWaitTimes(firstAttemptTime, secondAttemptTime),
				)
			},
			fn: func(attempt int) error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.ErrorIs(t, err, context.Canceled)
				assert.Equal(t, attempt, 1)
			},
		},
		{
			name: "success-context-canceled-in-meanwhile",
			ctx: func(t *testing.T) context.Context {
				return delayCancelContext(secondAttemptTime)
			},
			p: func(t *testing.T) retry.Policy {
				return retry.NewPolicy(
					retry.WithCount(3),
					retry.WithWaitTimes(firstAttemptTime, secondAttemptTime),
				)
			},
			fn: func(attempt int) error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.ErrorIs(t, err, context.Canceled)
				assert.Equal(t, true, attempt >= 2, "invalid attempt value", strconv.Itoa(attempt))
			},
		},
		{
			name: "success-context-deadline-exceeded",
			ctx: func(t *testing.T) context.Context {
				return deadlinedContext()
			},
			p: func(t *testing.T) retry.Policy {
				return retry.NewPolicy(
					retry.WithCount(3),
					retry.WithWaitTimes(firstAttemptTime, secondAttemptTime),
				)
			},
			fn: func(attempt int) error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Equal(t, attempt, 1)
			},
		},
		{
			name: "success-wait-times",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			p: func(t *testing.T) retry.Policy {
				return retry.NewPolicy(
					retry.WithCount(3),
					retry.WithWaitTimes(firstAttemptTime, secondAttemptTime),
					withTwoAttemptsWaitAssertion,
				)
			},
			fn: func(attempt int) error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.Error(t, err)
				assert.Equal(t, attempt, 4)
			},
		},
		{
			name: "success-waiter",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			p: func(t *testing.T) retry.Policy {
				return retry.NewPolicy(
					retry.WithCount(3),
					retry.WithWaiter(func(attempt int, err error) time.Duration {
						if attempt == 1 {
							return firstAttemptTime
						}

						if attempt == 2 {
							return secondAttemptTime
						}

						return 0
					}),
					withTwoAttemptsWaitAssertion,
				)
			},
			fn: func(attempt int) error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.Error(t, err)
				assert.Equal(t, attempt, 4)
			},
		},
		{
			name: "success-default-cancel",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			p: func(t *testing.T) retry.Policy {
				return retry.NewPolicy(
					retry.WithCount(1),
				)
			},
			fn: func(attempt int) error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, attempt int, err error) {
				assert.Error(t, err)
				assert.Equal(t, attempt, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zero := 0
			attempt := &zero
			p := tt.p(t)

			err := p.Execute(tt.ctx(t), func() error {
				*attempt++
				return tt.fn(*attempt)
			})

			tt.asserter(t, *attempt, err)
		})
	}
}

func delayCancelContext(delay time.Duration) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(delay, func() {
		cancel()
	})
	return ctx
}

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func deadlinedContext() context.Context {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	cancel()
	return ctx
}
