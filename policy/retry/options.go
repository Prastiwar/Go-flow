package retry

import (
	"time"
)

type Option func(*policy)

// CancelPredicate controls retrying continuation - return true to stop retrying and false to continue.
type CancelPredicate func(attempt int, err error) bool

// Waiter is a function called before next retry attempt will occur. Return time.Duration defining
// how much time program should wait before next retry execution.
type Waiter func(attempt int, err error) time.Duration

// WithCount configures how many times should retry be executed. If set to 0
// function will be executed one time and no retry will occur.
func WithCount(count int) Option {
	return func(rp *policy) {
		rp.count = count
	}
}

// WithWaitTimes configures Waiter that works as wait times enumerator for specified waitTimes.
// If there is difference between retry count and wait times - it will return the edge index on exceeded attempt.
func WithWaitTimes(waitTimes ...time.Duration) Option {
	l := len(waitTimes)
	if l == 0 {
		return WithWaiter(func(attempt int, err error) time.Duration {
			return 0
		})
	}

	return WithWaiter(func(attempt int, err error) time.Duration {
		attempt--
		if attempt >= l {
			return waitTimes[l-1]
		}
		if attempt < 0 {
			return waitTimes[0]
		}
		return waitTimes[attempt]
	})
}

// WithWaiter configures Waiter function to provide program wait time before retry execution.
func WithWaiter(waiter Waiter) Option {
	return func(rp *policy) {
		rp.waiter = waiter
	}
}

// WithCancelPredicate configures CancelPredicate to control retry continuation.
func WithCancelPredicate(handler CancelPredicate) Option {
	return func(rp *policy) {
		rp.cancel = handler
	}
}
