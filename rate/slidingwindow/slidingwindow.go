package slidingwindow

import (
	"errors"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
)

var (
	_ rate.Limiter = &limiter{}
	_ rate.Token   = &token{}
)

var (
	ErrInvalidInterval       = errors.New("interval must be non-zero positive duration value")
	ErrInvalidSegmentsTooLow = errors.New("segments value is too low")
	ErrInvalidLimitTooLow    = errors.New("maxEvents must be non-zero positive integer value")
)

const (
	zeroDuration = time.Duration(0)
)

type swOptions struct {
	Clock rate.Clock
}

type Option func(o *swOptions)

func newOptions(opts ...Option) *swOptions {
	o := &swOptions{}
	for _, opt := range opts {
		opt(o)
	}
	if o.Clock == nil {
		o.Clock = rate.SystemClock
	}
	return o
}

// WithClock returns an Option to modify clock used in time retrieval.
func WithClock(c rate.Clock) Option {
	return func(o *swOptions) {
		o.Clock = c
	}
}

// NewAlgorithm returns rate.LimiterAlgorithm with Sliding Window Counter implementation that's based on
// Fixed Window and Sliding Window algorithms to make a perfect balance between memory consumption and accuracy.
// The higher value of segments means better accuracy but also more memory consumption.
// The minimum value of segments is 2 and there is no check for maximum, however too big value can lead to inaccurate results
// so test the values for yourself to find the best accuracy to your usecase since it can vary depending on interval value.
// Segments are fixed-length time slots of value `interval/segments`. Each has an internal counter.
// The last segment is the current one which is incrementing, after segment duration, it's shifted to the back of 1 segment slot
// so if the segment duration passes, the last has 0 value, and the oldest is replaced by a counter at 1 index e.g [1, 2, 3, 4, 5] => [2, 3, 4, 5, 0].
// It's recommended to use for frequent events in relatively small interval like web API rate limitation.
func NewAlgorithm(maxEvents uint32, interval time.Duration, segments uint32, options ...Option) (rate.LimiterAlgorithm, error) {
	if maxEvents <= 0 {
		return nil, ErrInvalidLimitTooLow
	}
	if interval <= zeroDuration {
		return nil, ErrInvalidInterval
	}
	if segments <= 1 {
		return nil, ErrInvalidSegmentsTooLow
	}

	o := newOptions(options...)

	return func() rate.Limiter {
		return newLimiter(maxEvents, interval, segments, o.Clock)
	}, nil
}

type windowState struct {
	maxEvents uint32
	interval  time.Duration
	counters  []uint32
	startTime time.Time
}

func (l *windowState) advance(now time.Time) {
	segments := len(l.counters)
	lastSegmentIndex := segments - 1
	segmentInterval := l.interval / time.Duration(segments)

	// if request time passed whole window, then counters should be reset
	if now.Sub(l.startTime) >= l.interval {
		l.startTime = now
		for i := range l.counters {
			l.counters[i] = 0
		}
	}

	for now.Sub(l.startTime) >= segmentInterval {
		l.startTime = l.startTime.Add(segmentInterval)
		// swap elements from start in left direction
		// so first element in fact disappears and the last will be 0
		// e.g [1, 2, 3, 4, 5] => [2, 3, 4, 5, 0]
		if lastSegmentIndex%2 == 0 {
			for i := 0; i < lastSegmentIndex; i += 2 {
				j := i + 1
				l.counters[i], l.counters[j] = l.counters[j], l.counters[j+1]
			}
		} else {
			for i := 0; i < lastSegmentIndex; i++ {
				l.counters[i] = l.counters[i+1]
			}
		}
		l.counters[lastSegmentIndex] = 0
	}
}

func (l *windowState) Count(now time.Time) uint32 {
	l.advance(now)

	total := uint32(0)
	for _, count := range l.counters {
		total += count
	}

	return total
}

// Available returns amount of tokens that could be consumed.
func (l *windowState) Available(now time.Time) uint64 {
	total := l.Count(now)
	return l.clamp(int32(l.maxEvents) - int32(total))
}

func (l *windowState) clamp(v int32) uint64 {
	if v <= 0 {
		return 0
	}
	return uint64(v)
}

// Incr increments the last counter which is always current segment.
func (l *windowState) Incr() {
	l.counters[len(l.counters)-1]++
}

type limiter struct {
	state *windowState
	clock rate.Clock
}

func newLimiter(maxEvents uint32, interval time.Duration, segments uint32, clock rate.Clock) *limiter {
	return &limiter{
		state: &windowState{
			maxEvents: maxEvents,
			interval:  interval,
			counters:  make([]uint32, segments),
			startTime: clock.Now(),
		},
		clock: clock,
	}
}

func (l *limiter) Limit() uint64 {
	return uint64(l.state.maxEvents)
}

func (l *limiter) Take() rate.Token {
	return newToken(l.state, l.clock)
}

func (l *limiter) Tokens() uint64 {
	now := l.clock.Now()
	return l.state.Available(now)
}

type token struct {
	state *windowState
	clock rate.Clock
}

func newToken(l *windowState, clock rate.Clock) rate.Token {
	return &token{state: l, clock: clock}
}

func (t *token) Allow() bool {
	return t.state.Available(t.clock.Now()) > 0
}

func (t *token) ResetsAt() time.Time {
	segmentInterval := t.state.interval / time.Duration(len(t.state.counters))
	now := t.clock.Now()

	if t.Allow() {
		return now
	}

	for i, v := range t.state.counters {
		if v > 0 {
			targetDur := segmentInterval * (time.Duration(i + 1))
			actualDur := now.Sub(t.state.startTime)
			now = now.Add(targetDur - actualDur)
			break
		}
	}
	return now
}

func (t *token) Use() error {
	if ok := t.Allow(); !ok {
		return rate.ErrRateLimitExceeded
	}
	t.state.Incr()
	return nil
}
