package slidingwindow_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/rate/slidingwindow"
	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

func TestSlidingToken(t *testing.T) {
	// Arrange
	clock, setTimer := mocks.NewMutableClock()

	alg, err := slidingwindow.NewAlgorithm(5, 15*time.Second, 3, slidingwindow.WithClock(clock))
	if err != nil {
		panic(err)
	}
	limiter := alg()
	token := limiter.Take()

	// Act & Assert
	assert.Equal(t, uint64(5), limiter.Limit(), "pre limit")
	assert.Equal(t, uint64(5), limiter.Tokens(), "pre tokens")

	i := 0
	for token.Use() == nil {
		i++
	}
	assert.ErrorIs(t, token.Use(), rate.ErrRateLimitExceeded, "post limit")
	assert.Equal(t, uint64(5), limiter.Limit(), "post limit")
	assert.Equal(t, uint64(0), limiter.Tokens(), "post tokens")
	assert.Equal(t, 5, i, "use count")

	setTimer(clock.Now().Add(30 * time.Second))

	assert.Equal(t, uint64(5), limiter.Tokens(), "reset tokens")
}

func TestSlidingWindowCounter(t *testing.T) {
	tests := []struct {
		name        string
		maxEvents   uint32
		interval    time.Duration
		segments    uint32
		timeMod     func(i int, resetsAt time.Time, clock rate.Clock) time.Time
		expectation []time.Duration
	}{
		{
			name:      "5-per-15s",
			maxEvents: 5,
			interval:  15 * time.Second,
			segments:  3,
			expectation: []time.Duration{
				/*  [0] */ 0 * time.Second,
				/*  [1] */ 14 * time.Second,
				/*  [2] */ 14 * time.Second,
				/*  [3] */ 14 * time.Second,
				/*  [4] */ 14 * time.Second,
				/*  [5] */ 15 * time.Second,
				/*  [6] */ 25 * time.Second,
				/*  [7] */ 25 * time.Second,
				/*  [8] */ 25 * time.Second,
				/*  [9] */ 25 * time.Second,
				/* [10] */ 30 * time.Second,
				/* [11] */ 40 * time.Second,
				/* [12] */ 40 * time.Second,
				/* [13] */ 40 * time.Second,
				/* [14] */ 70 * time.Second,
			},
			timeMod: func(i int, resetsAt time.Time, clock rate.Clock) time.Time {
				if i == 0 {
					return clock.Now().Add(14 * time.Second)
				} else if i == 13 {
					return clock.Now().Add(30 * time.Second)
				}
				return resetsAt
			},
		},
		{
			name:      "4-per-10s",
			maxEvents: 4,
			interval:  10 * time.Second,
			segments:  4,
			expectation: []time.Duration{
				/*  [0] */ 0 * time.Second,
				/*  [1] */ 1 * time.Second,
				/*  [2] */ 2 * time.Second,
				/*  [3] */ 3 * time.Second,
				/*  [4] */ 11 * time.Second,
				/*  [5] */ 12 * time.Second,
				/*  [6] */ 13 * time.Second,
				/*  [7] */ 14 * time.Second,
				/*  [8] */ 21 * time.Second,
				/*  [9] */ 22 * time.Second,
				/* [10] */ 23 * time.Second,
				/* [11] */ 24 * time.Second,
				/* [12] */ 31 * time.Second,
				/* [13] */ 32 * time.Second,
				/* [14] */ 33 * time.Second,
			},
			timeMod: func(i int, resetsAt time.Time, clock rate.Clock) time.Time {
				return resetsAt.Add(time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			clock, setTimer := mocks.NewMutableClock()

			alg, err := slidingwindow.NewAlgorithm(tt.maxEvents, tt.interval, tt.segments, slidingwindow.WithClock(clock))
			if err != nil {
				panic(err)
			}
			limiter := alg()
			times := make([]time.Duration, len(tt.expectation))

			// Act
			start := clock.Now()
			for i := 0; i < len(tt.expectation); i++ {
				t := limiter.Take()
				if err := t.Use(); err != nil {
					panic(err)
				}

				times[i] = clock.Now().Sub(start)

				if tt.timeMod != nil {
					targetTime := tt.timeMod(i, t.ResetsAt(), clock)
					setTimer(targetTime)
					continue
				}

				resetsAt := t.ResetsAt()
				setTimer(resetsAt)
			}

			// Assert
			delta := time.Second / 2
			for i, val := range times {
				assert.Approximately(t, tt.expectation[i], val, delta, "index at ["+strconv.Itoa(i)+"]")
			}
		})
	}
}

func TestSlidingWindowCounterConstructor(t *testing.T) {
	tests := []struct {
		name      string
		maxEvents uint32
		interval  time.Duration
		segments  uint32
		options   []slidingwindow.Option
		assertErr assert.ErrorFunc
	}{
		{
			name:      "valid",
			maxEvents: 10,
			interval:  time.Minute,
			segments:  15,
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:      "invalid-maxEvents",
			maxEvents: 0,
			assertErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, slidingwindow.ErrInvalidLimitTooLow)
			},
		},
		{
			name:      "invalid-interval",
			maxEvents: 1,
			interval:  time.Duration(-1),
			segments:  2,
			assertErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, slidingwindow.ErrInvalidInterval)
			},
		},
		{
			name:      "invalid-segments-low",
			maxEvents: 1,
			interval:  time.Second,
			segments:  1,
			assertErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, slidingwindow.ErrInvalidSegmentsTooLow)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := slidingwindow.NewAlgorithm(tt.maxEvents, tt.interval, tt.segments, tt.options...)
			tt.assertErr(t, err)
		})
	}
}
