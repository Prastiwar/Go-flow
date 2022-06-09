package assert

import (
	"testing"
)

type Counter struct {
	expectedCount int
	counter       int
}

func Count(n int) *Counter {
	return &Counter{expectedCount: n}
}

func (c *Counter) Inc() {
	c.counter++
}

func (c *Counter) Assert(t *testing.T) {
	if c.expectedCount != c.counter {
		t.Errorf("expected call count: '%v' actual: '%v'", c.expectedCount, c.counter)
	}
}
