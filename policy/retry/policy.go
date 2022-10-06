package retry

import "time"

type Policy struct {
	count  int
	waiter Waiter
	cancel CancelPredicate
}

func NewPolicy(opts ...Option) *Policy {
	p := &Policy{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Policy) Execute(fn func() error) error {
	var err error
	attempts := p.count + 1
	cancel := p.cancel

	if cancel == nil {
		cancel = func(attempt int, err error) bool {
			return false
		}
	}

	for i := 1; i <= attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}

		stop := cancel(i, err)
		if stop {
			break
		}

		p.wait(i, err)
	}

	return err
}

func (p *Policy) wait(attempt int, err error) {
	if p.waiter == nil {
		return
	}

	if attempt == p.count+1 {
		// there is no further attempts so no need to wait
		return
	}

	waitTime := p.waiter(attempt, err)

	time.Sleep(waitTime)
}
