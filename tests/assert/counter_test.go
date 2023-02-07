package assert_test

import (
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestCounterAutoAssert(t *testing.T) {
	// there is no way to assert if t.Cleanup was called or to remove
	// failed status, so need to look for coverage if this passes correctly
	assert.Count(t, 0)
}

func TestCounterAssert(t *testing.T) {
	tests := []struct {
		name string
		c    func(t *testing.T) *assert.Counter
		fail bool
	}{
		{
			name: "success-zero-assertion",
			c: func(t *testing.T) *assert.Counter {
				return assert.Count(t, 0)
			},
			fail: false,
		},
		{
			name: "success-assertion",
			c: func(t *testing.T) *assert.Counter {
				c := assert.Count(t, 1)
				c.Inc()
				return c
			},
			fail: false,
		},
		{
			name: "fail-assertion",
			c: func(t *testing.T) *assert.Counter {
				return assert.Count(t, 1)
			},
			fail: true,
		},
		{
			name: "at-least-assertion-with-more",
			c: func(t *testing.T) *assert.Counter {
				c := assert.Count(t, 1).AtLeast()
				c.Inc()
				c.Inc()
				return c
			},
			fail: false,
		},
		{
			name: "at-least-assertion-with-less",
			c: func(t *testing.T) *assert.Counter {
				return assert.Count(t, 1).AtLeast()
			},
			fail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			tt.c(test).Assert(test)

			assert.Equal(t, tt.fail, test.Failed())
		})
	}
}
