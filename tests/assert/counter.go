package assert

import (
	"fmt"
	"testing"
)

type Counter struct {
	expectedCount int
	counter       int
	atLeast       bool
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

// Inc marks a single call to the function.
func (c *Counter) Inc() {
	c.counter++
}

// Assert checks for count expectation.
func (c *Counter) Assert(t *testing.T, prefixes ...string) {
	t.Helper()
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
