package rate

import (
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestSystemClock(t *testing.T) {
	clock := SystemClock

	assert.Equal(t, time.Now(), clock.Now())
}
