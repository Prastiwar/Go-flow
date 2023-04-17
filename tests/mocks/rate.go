package mocks

import (
	"context"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

var (
	_ rate.LimiterStore            = LimiterStoreMock{}
	_ rate.Limiter                 = LimiterMock{}
	_ rate.BurstLimiter            = BurstLimiterMock{}
	_ rate.ReservationLimiter      = ReservationLimiterMock{}
	_ rate.BurstReservationLimiter = BurstReservationLimiterMock{}
	_ rate.CancellableToken        = TokenMock{}
	_ rate.Clock                   = MockClock{}
)

type LimiterStoreMock struct {
	OnLimit func(ctx context.Context, key string) (rate.Limiter, error)
}

func (m LimiterStoreMock) Limit(ctx context.Context, key string) (rate.Limiter, error) {
	assert.ExpectCall(m.OnLimit)
	return m.OnLimit(ctx, key)
}

type LimiterMock struct {
	OnLimit  func() uint64
	OnTokens func(ctx context.Context) (uint64, error)
	OnTake   func(ctx context.Context) (rate.Token, error)
}

func (m LimiterMock) Limit() uint64 {
	assert.ExpectCall(m.OnLimit)
	return m.OnLimit()
}

func (m LimiterMock) Take(ctx context.Context) (rate.Token, error) {
	assert.ExpectCall(m.OnTake)
	return m.OnTake(ctx)
}

func (m LimiterMock) Tokens(ctx context.Context) (uint64, error) {
	assert.ExpectCall(m.OnTokens)
	return m.OnTokens(ctx)
}

type BurstLimiterMock struct {
	rate.Limiter

	OnBurst func() uint64
	OnTakeN func(ctx context.Context, n uint64) (rate.Token, error)
}

func (m BurstLimiterMock) TakeN(ctx context.Context, n uint64) (rate.Token, error) {
	assert.ExpectCall(m.OnTakeN)
	return m.OnTakeN(ctx, n)
}

func (m BurstLimiterMock) Burst() uint64 {
	assert.ExpectCall(m.OnBurst)
	return m.OnBurst()
}

type ReservationLimiterMock struct {
	rate.Limiter

	OnReserve func(ctx context.Context) (rate.CancellableToken, error)
}

func (m ReservationLimiterMock) Reserve(ctx context.Context) (rate.CancellableToken, error) {
	assert.ExpectCall(m.OnReserve)
	return m.OnReserve(ctx)
}

type BurstReservationLimiterMock struct {
	rate.BurstLimiter

	OnReserve  func(ctx context.Context) (rate.CancellableToken, error)
	OnReserveN func(ctx context.Context, n uint64) (rate.CancellableToken, error)
}

func (m BurstReservationLimiterMock) Reserve(ctx context.Context) (rate.CancellableToken, error) {
	assert.ExpectCall(m.OnReserve)
	return m.OnReserve(ctx)
}

func (m BurstReservationLimiterMock) ReserveN(ctx context.Context, n uint64) (rate.CancellableToken, error) {
	assert.ExpectCall(m.OnReserveN)
	return m.OnReserveN(ctx, n)
}

type TokenMock struct {
	OnAllow    func() bool
	OnResetsAt func() time.Time
	OnUse      func() error
	OnCancel   func()
	OnContext  func() context.Context
}

func (m TokenMock) Allow() bool {
	assert.ExpectCall(m.OnAllow)
	return m.OnAllow()
}

func (m TokenMock) ResetsAt() time.Time {
	assert.ExpectCall(m.OnResetsAt)
	return m.OnResetsAt()
}

func (m TokenMock) Use() error {
	assert.ExpectCall(m.OnUse)
	return m.OnUse()
}

func (m TokenMock) Cancel() {
	assert.ExpectCall(m.OnCancel)
	m.OnCancel()
}

func (m TokenMock) Context() context.Context {
	assert.ExpectCall(m.OnContext)
	return m.OnContext()
}

type MockClock struct {
	NowFunc func() time.Time
}

func (c MockClock) Now() time.Time {
	return c.NowFunc()
}

// NewMutableClock returns MockClock and function to set time of the clock at runtime.
func NewMutableClock() (rate.Clock, func(time.Time)) {
	now := time.Now()
	timer := &now

	mockClock := MockClock{
		NowFunc: func() time.Time {
			return *timer
		},
	}

	return mockClock, func(t time.Time) { *timer = t }
}
