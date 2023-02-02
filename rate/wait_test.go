package rate

import (
	"context"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

var (
	defaultDelta = 50 * time.Millisecond
)

type limiterDummy struct {
	resetAt time.Time
}

func TestWait(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		resetAt   time.Time
		assertErr assert.ResultErrorFunc[time.Duration]
	}{
		{
			name:    "success-non-waiting",
			ctx:     context.Background(),
			resetAt: time.Now(),
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.NilError(t, err)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name:    "success-wait-half-second",
			ctx:     context.Background(),
			resetAt: time.Now().Add(time.Second / 2),
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.NilError(t, err)
				assert.Approximately(t, time.Second/2, result, defaultDelta)
			},
		},
		{
			name:    "failure-wait-canceled",
			ctx:     canceledContext(),
			resetAt: time.Now().Add(time.Second / 2),
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorIs(t, err, context.Canceled)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name:    "failure-wait-deadline-exceeded",
			ctx:     deadlinedContext(),
			resetAt: time.Now().Add(time.Second / 2),
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			err := Wait(tt.ctx, tt.resetAt)

			tt.assertErr(t, time.Since(startTime), err)
		})
	}
}

// func TestTakeAndWait(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		ctx       context.Context
// 		limiter   Limiter
// 		assertErr assert.ResultErrorFunc[time.Duration]
// 	}{
// 		{
// 			name: "success-non-waiting",
// 			ctx:  context.Background(),
// 			limiter: LimiterMock{
// 				time.Now(),
// 			},
// 			assertErr: func(t *testing.T, result time.Duration, err error) {
// 				assert.NilError(t, err)
// 				assert.Approximately(t, time.Duration(0), result, defaultDelta)
// 			},
// 		},
// 		{
// 			name: "success-wait-half-second",
// 			ctx:  context.Background(),
// 			limiter: LimiterMock{
// 				time.Now().Add(time.Second / 2),
// 			},
// 			assertErr: func(t *testing.T, result time.Duration, err error) {
// 				assert.NilError(t, err)
// 				assert.Approximately(t, time.Second/2, result, defaultDelta)
// 			},
// 		},
// 		{
// 			name: "failure-wait-canceled",
// 			ctx:  canceledContext(),
// 			limiter: LimiterMock{
// 				time.Now().Add(time.Second / 2),
// 			},
// 			assertErr: func(t *testing.T, result time.Duration, err error) {
// 				assert.ErrorIs(t, err, context.Canceled)
// 				assert.Approximately(t, time.Duration(0), result, defaultDelta)
// 			},
// 		},
// 		{
// 			name: "failure-wait-deadline-exceeded",
// 			ctx:  deadlinedContext(),
// 			limiter: LimiterMock{
// 				time.Now().Add(time.Second / 2),
// 			},
// 			assertErr: func(t *testing.T, result time.Duration, err error) {
// 				assert.ErrorIs(t, err, context.DeadlineExceeded)
// 				assert.Approximately(t, time.Duration(0), result, defaultDelta)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			startTime := time.Now()
// 			err := TakeAndWait(tt.ctx, tt.limiter)

// 			tt.assertErr(t, time.Since(startTime), err)
// 		})
// 	}
// }

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func deadlinedContext() context.Context {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	cancel()
	return ctx
}
