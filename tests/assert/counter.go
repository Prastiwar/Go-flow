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
	Equal(t, c.expectedCount, c.counter)
}
