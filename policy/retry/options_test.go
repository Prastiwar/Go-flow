package retry

import (
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestInvalidWithWaitTimes(t *testing.T) {
	waitTimes := []time.Duration{}
	p := NewPolicy(WithWaitTimes(waitTimes...))

	assert.NotNil(t, p.waiter, "waiter expectation failed")

	actualDur := p.waiter(-5, nil)
	assert.Equal(t, time.Duration(0), actualDur)

	actualDur = p.waiter(len(waitTimes)+1, nil)
	assert.Equal(t, time.Duration(0), actualDur)
}

func TestWithWaitTimes(t *testing.T) {
	waitTimes := []time.Duration{
		time.Second, 2 * time.Second, 3 * time.Second,
	}
	p := NewPolicy(WithWaitTimes(waitTimes...))

	assert.NotNil(t, p.waiter, "waiter expectation failed")

	for i, dur := range waitTimes {
		actualDur := p.waiter(i+1, nil)
		assert.Equal(t, dur, actualDur)
	}

	actualDur := p.waiter(-5, nil)
	assert.Equal(t, waitTimes[0], actualDur)

	actualDur = p.waiter(len(waitTimes)+1, nil)
	assert.Equal(t, waitTimes[len(waitTimes)-1], actualDur)
}
