package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestPolicy_Execute(t *testing.T) {
	const (
		firstAttemptTime  = 100 * time.Millisecond
		secondAttemptTime = 200 * time.Millisecond
		deltaDivider      = 2
	)

	var attempt int
	var timeStart time.Time

	withTwoAttemptsWaitAssertion := WithCancelPredicate(func(attempt int, err error) bool {
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
		p        func(t *testing.T) *Policy
		fn       func() error
		asserter assert.ErrorFunc
	}{
		{
			name: "success-after-single-retry",
			p: func(t *testing.T) *Policy {
				c := assert.Count(t, 1, "failed retry call")
				return NewPolicy(
					WithCount(1),
					WithCancelPredicate(func(attempt int, err error) bool {
						c.Inc()
						return false
					}),
				)
			},
			fn: func() error {
				if attempt == 0 {
					attempt++
					return errors.New("invalid")
				}
				return nil
			},
			asserter: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name: "success-3-retries-cancel",
			p: func(t *testing.T) *Policy {
				c := assert.Count(t, 3, "failed retry call")
				return NewPolicy(
					WithCount(5),
					WithCancelPredicate(func(attempt int, err error) bool {
						c.Inc()
						return attempt == 3
					}),
				)
			},
			fn: func() error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success-no-cancel-no-error",
			p: func(t *testing.T) *Policy {
				return NewPolicy(
					WithCount(5),
				)
			},
			fn: func() error {
				if attempt == 1 {
					attempt++
					return errors.New("invalid")
				}
				return nil
			},
			asserter: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name: "success-wait-times",
			p: func(t *testing.T) *Policy {
				return NewPolicy(
					WithCount(3),
					WithWaitTimes(firstAttemptTime, secondAttemptTime),
					withTwoAttemptsWaitAssertion,
				)
			},
			fn: func() error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success-waiter",
			p: func(t *testing.T) *Policy {
				return NewPolicy(
					WithCount(3),
					WithWaiter(func(attempt int, err error) time.Duration {
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
			fn: func() error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success-default-cancel",
			p: func(t *testing.T) *Policy {
				return NewPolicy(
					WithCount(1),
				)
			},
			fn: func() error {
				return errors.New("invalid")
			},
			asserter: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempt = 0
			p := tt.p(t)

			err := p.Execute(tt.fn)

			tt.asserter(t, err)
		})
	}
}
