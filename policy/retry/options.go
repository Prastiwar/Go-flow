package retry

import (
	"time"
)

type Option func(*Policy)

// return true to stop retrying and false to continue retrying
type CancelPredicate func(attempt int, err error) bool

type Waiter func(attempt int, err error) time.Duration

func WithCount(count int) Option {
	return func(rp *Policy) {
		rp.count = count
	}
}

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

func WithWaiter(waiter Waiter) Option {
	return func(rp *Policy) {
		rp.waiter = waiter
	}
}

func WithCancelPredicate(handler CancelPredicate) Option {
	return func(rp *Policy) {
		rp.cancel = handler
	}
}
