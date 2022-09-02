package assert

import (
	"testing"
)

func TestCounterAssert(t *testing.T) {
	tests := []struct {
		name string
		c    func(t *testing.T) *Counter
		fail bool
	}{
		{
			name: "success-zero-assertion",
			c: func(t *testing.T) *Counter {
				return Count(0)
			},
			fail: false,
		},
		{
			name: "success-assertion",
			c: func(t *testing.T) *Counter {
				c := Count(1)
				c.Inc()
				return c
			},
			fail: false,
		},
		{
			name: "fail-assertion",
			c: func(t *testing.T) *Counter {
				return Count(1)
			},
			fail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testing.T{}
			tt.c(t).Assert(test)
			Equal(t, tt.fail, test.Failed())
		})
	}
}
