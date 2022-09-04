package assert

import (
	"testing"
)

type Counter struct {
	expectedCount int
	counter       int
}

// Count returns a Counter with expected n count to be asserted. Use Inc() to mark a call.
// Assert() will be called on test cleanup, so it's not necessary to call it manually but recommended
// due to convention to not hide any test side effects.
func Count(t *testing.T, n int) *Counter {
	c := &Counter{expectedCount: n}
	t.Cleanup(func() {
		c.Assert(t)
	})

	return c
}

// Inc marks a single call to the function.
func (c *Counter) Inc() {
	c.counter++
}

// Assert checks for count expectation.
func (c *Counter) Assert(t *testing.T) {
	if c.expectedCount != c.counter {
		t.Errorf("expected call count: '%v' actual: '%v'", c.expectedCount, c.counter)
	}
}
