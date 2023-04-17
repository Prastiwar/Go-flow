package assert_test

import (
	"sync"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestCounterAutoAssert(t *testing.T) {
	// there is no way to assert if t.Cleanup was called or to remove
	// failed status, so need to look for coverage if this passes correctly
	assert.Count(t, 0)
}

// TestCounterAlwaysOneAssert relies on code coverage and will not fail is multiple asserts are called.
// There is no way for any kind of assertions on testing.T.
func TestCounterAlwaysOneAssert(t *testing.T) {
	const count = 10

	testT := &testing.T{}
	c := assert.Count(testT, 1)
	wg := &sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			c.Assert(testT)
		}()
	}
	wg.Wait()

	assert.Equal(t, true, testT.Failed())
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
