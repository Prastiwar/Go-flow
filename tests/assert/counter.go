package assert

import (
	"fmt"
	"sync"
	"testing"
)

type Counter struct {
	expectedCount int
	counter       int
	atLeast       bool

	done bool
	m    sync.Mutex
}

// Count returns a Counter with expected n count to be asserted. Use Inc() to mark a call.
// Assert() will be called on test cleanup, so it's not necessary to call it manually but recommended
// due to convention to not hide any test side effects.
func Count(t *testing.T, n int, prefixes ...string) *Counter {
	t.Helper()
	c := &Counter{expectedCount: n}
	t.Cleanup(func() {
		t.Helper()
		c.Assert(t, prefixes...)
	})

	return c
}

// AtLeast marks counter to assert for at least count instead exact count and returns self.
func (c *Counter) AtLeast() *Counter {
	c.atLeast = true
	return c
}

// Inc marks a single call to the function. This is thread safe operation.
func (c *Counter) Inc() {
	c.m.Lock()
	defer c.m.Unlock()

	c.counter++
}

// Assert checks for count expectation.
func (c *Counter) Assert(t *testing.T, prefixes ...string) {
	c.m.Lock()
	defer c.m.Unlock()

	t.Helper()
	if c.done {
		return
	}
	c.done = true

	if c.atLeast {
		if c.expectedCount > c.counter {
			errorf(t, fmt.Sprintf("expected call count at least: '%v', actual: '%v'", c.expectedCount, c.counter), prefixes...)
		}
		return
	}
	if c.expectedCount != c.counter {
		errorf(t, fmt.Sprintf("expected call count: '%v', actual: '%v'", c.expectedCount, c.counter), prefixes...)
	}
}
