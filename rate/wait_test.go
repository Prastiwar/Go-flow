package rate_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

var (
	defaultDelta = 50 * time.Millisecond
)

func TestWait(t *testing.T) {
	tests := []struct {
		name      string
		ctx       func(t *testing.T) context.Context
		resetAt   time.Time
		assertErr assert.ResultErrorFunc[time.Duration]
	}{
		{
			name: "success-non-waiting",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			resetAt: time.Now(),
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.NilError(t, err)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "success-wait-half-second",
			ctx: func(t *testing.T) context.Context {
				return delayCancelContext(time.Second)
			},
			resetAt: time.Now().Add(time.Second / 2),
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.NilError(t, err)
				assert.Approximately(t, time.Second/2, result, defaultDelta)
			},
		},
		{
			name: "failure-wait-canceled",
			ctx: func(t *testing.T) context.Context {
				return canceledContext()
			},
			resetAt: time.Now().Add(time.Second / 2),
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorIs(t, err, context.Canceled)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "failure-wait-deadline-exceeded",
			ctx: func(t *testing.T) context.Context {
				return deadlinedContext()
			},
			resetAt: time.Now().Add(time.Second / 2),
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "failure-wait-parent-canceled-first",
			ctx: func(t *testing.T) context.Context {
				return delayCancelContext(time.Second / 3)
			},
			resetAt: time.Now().Add(time.Second),
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorIs(t, err, context.Canceled)
				assert.Approximately(t, time.Duration(time.Second/3), result, defaultDelta)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()

			err := rate.Wait(tt.ctx(t), tt.resetAt)

			tt.assertErr(t, time.Since(startTime), err)
		})
	}
}

func TestConsumeAndWait(t *testing.T) {
	tests := []struct {
		name      string
		ctx       func(t *testing.T) context.Context
		limiter   rate.Limiter
		assertErr assert.ResultErrorFunc[time.Duration]
	}{
		{
			name: "success-non-waiting",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			limiter: mocks.LimiterMock{
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return nil
						},
						OnResetsAt: func() time.Time {
							return time.Now()
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.NilError(t, err)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "success-wait-half-second",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			limiter: mocks.LimiterMock{
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return rate.ErrRateLimitExceeded
						},
						OnResetsAt: func() time.Time {
							return time.Now().Add(time.Second / 2)
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.NilError(t, err)
				assert.Approximately(t, time.Second/2, result, defaultDelta)
			},
		},
		{
			name: "failure-wait-canceled",
			ctx: func(t *testing.T) context.Context {
				return canceledContext()
			},
			limiter: mocks.LimiterMock{
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return rate.ErrRateLimitExceeded
						},
						OnResetsAt: func() time.Time {
							return time.Now().Add(time.Second)
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorIs(t, err, context.Canceled)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "failure-wait-deadline-exceeded",
			ctx: func(t *testing.T) context.Context {
				return deadlinedContext()
			},
			limiter: mocks.LimiterMock{
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return rate.ErrRateLimitExceeded
						},
						OnResetsAt: func() time.Time {
							return time.Now().Add(time.Second)
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "failure-use-error",
			ctx: func(t *testing.T) context.Context {
				return deadlinedContext()
			},
			limiter: mocks.LimiterMock{
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return errors.New("use-error")
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorWith(t, err, "use-error")
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "failure-take-error",
			ctx: func(t *testing.T) context.Context {
				return deadlinedContext()
			},
			limiter: mocks.LimiterMock{
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return nil, errors.New("take-error")
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorWith(t, err, "take-error")
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			err := rate.ConsumeAndWait(tt.ctx(t), tt.limiter)

			tt.assertErr(t, time.Since(startTime), err)
		})
	}
}

func TestConsumeNAndWait(t *testing.T) {
	tests := []struct {
		name      string
		ctx       func(t *testing.T) context.Context
		limiter   rate.BurstLimiter
		n         uint64
		assertErr assert.ResultErrorFunc[time.Duration]
	}{
		{
			name: "success-non-waiting",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			limiter: mocks.BurstLimiterMock{
				OnTakeN: func(ctx context.Context, n uint64) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return nil
						},
						OnResetsAt: func() time.Time {
							return time.Now()
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.NilError(t, err)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "success-wait-half-second",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			limiter: mocks.BurstLimiterMock{
				OnTakeN: func(ctx context.Context, n uint64) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return rate.ErrRateLimitExceeded
						},
						OnResetsAt: func() time.Time {
							return time.Now().Add(time.Second / 2)
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.NilError(t, err)
				assert.Approximately(t, time.Second/2, result, defaultDelta)
			},
		},
		{
			name: "failure-wait-canceled",
			ctx: func(t *testing.T) context.Context {
				return canceledContext()
			},
			limiter: mocks.BurstLimiterMock{
				OnTakeN: func(ctx context.Context, n uint64) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return rate.ErrRateLimitExceeded
						},
						OnResetsAt: func() time.Time {
							return time.Now().Add(time.Second)
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorIs(t, err, context.Canceled)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "failure-wait-deadline-exceeded",
			ctx: func(t *testing.T) context.Context {
				return deadlinedContext()
			},
			limiter: mocks.BurstLimiterMock{
				OnTakeN: func(ctx context.Context, n uint64) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return rate.ErrRateLimitExceeded
						},
						OnResetsAt: func() time.Time {
							return time.Now().Add(time.Second)
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "failure-wait-deadline-exceeded",
			ctx: func(t *testing.T) context.Context {
				return deadlinedContext()
			},
			limiter: mocks.BurstLimiterMock{
				OnTakeN: func(ctx context.Context, n uint64) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse: func() error {
							return errors.New("use-error")
						},
					}, nil
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorWith(t, err, "use-error")
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
		{
			name: "failure-take-error",
			ctx: func(t *testing.T) context.Context {
				return deadlinedContext()
			},
			limiter: mocks.BurstLimiterMock{
				OnTakeN: func(ctx context.Context, n uint64) (rate.Token, error) {
					return nil, errors.New("take-error")
				},
			},
			assertErr: func(t *testing.T, result time.Duration, err error) {
				assert.ErrorWith(t, err, "take-error")
				assert.Approximately(t, time.Duration(0), result, defaultDelta)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()

			err := rate.ConsumeNAndWait(tt.ctx(t), tt.limiter, tt.n)

			tt.assertErr(t, time.Since(startTime), err)
		})
	}
}

func delayCancelContext(delay time.Duration) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(delay, func() {
		cancel()
	})
	return ctx
}

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
