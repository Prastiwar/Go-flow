package assert

import (
	"testing"
)

func TestCounterAutoAssert(t *testing.T) {
	// there is no way to assert if t.Cleanup was called or to remove
	// failed status, so need to look for coverage if this passes correctly
	Count(t, 0)
}

func TestCounterAssert(t *testing.T) {
	tests := []struct {
		name string
		c    func(t *testing.T) *Counter
		fail bool
	}{
		{
			name: "success-zero-assertion",
			c: func(t *testing.T) *Counter {
				return Count(t, 0)
			},
			fail: false,
		},
		{
			name: "success-assertion",
			c: func(t *testing.T) *Counter {
				c := Count(t, 1)
				c.Inc()
				return c
			},
			fail: false,
		},
		{
			name: "fail-assertion",
			c: func(t *testing.T) *Counter {
				return Count(t, 1)
			},
			fail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}

			tt.c(test).Assert(test)

			Equal(t, tt.fail, test.Failed())
		})
	}
}
