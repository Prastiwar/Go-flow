package rate

import (
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestSystemClock(t *testing.T) {
	clock := SystemClock

	assert.ApproximatelyTime(t, time.Now(), clock.Now(), 100*time.Millisecond)
}
