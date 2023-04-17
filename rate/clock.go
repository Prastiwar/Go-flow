package rate

import "time"

var (
	// SystemClock implements Clock interface and returns the current local time with time.Now().
	SystemClock Clock = &systemClock{}
)

// Clock is type which defines Now() to get actual time.Time value.
// Most of the time SystemClock can be used to fulfill the interface.
type Clock interface {
	Now() time.Time
}

type systemClock struct{}

// Now returns the current local time.
func (c *systemClock) Now() time.Time {
	return time.Now()
}
