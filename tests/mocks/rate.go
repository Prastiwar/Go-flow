package mocks

import (
	"time"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

var (
	_ rate.LimiterStore     = LimiterStoreMock{}
	_ rate.Limiter          = LimiterMock{}
	_ rate.BurstLimiter     = BurstLimiterMock{}
	_ rate.CancellableToken = TokenMock{}
)

type LimiterStoreMock struct {
	OnLimit func(key string) rate.Limiter
}

func (m LimiterStoreMock) Limit(key string) rate.Limiter {
	assert.ExpectCall(m.OnLimit)
	return m.OnLimit(key)
}

type LimiterMock struct {
	OnLimit  func() uint64
	OnTokens func() uint64
	OnTake   func() rate.Token
}

func (m LimiterMock) Limit() uint64 {
	assert.ExpectCall(m.OnLimit)
	return m.OnLimit()
}

func (m LimiterMock) Take() rate.Token {
	assert.ExpectCall(m.OnTake)
	return m.OnTake()
}

func (m LimiterMock) Tokens() uint64 {
	assert.ExpectCall(m.OnTokens)
	return m.OnTokens()
}

type BurstLimiterMock struct {
	rate.Limiter

	OnBurst func() uint64
	OnTakeN func(n uint64) rate.Token
}

func (m BurstLimiterMock) TakeN(n uint64) rate.Token {
	assert.ExpectCall(m.OnTakeN)
	return m.OnTakeN(n)
}

func (m BurstLimiterMock) Burst() uint64 {
	assert.ExpectCall(m.OnBurst)
	return m.OnBurst()
}

type TokenMock struct {
	OnAllow    func() bool
	OnResetsAt func() time.Time
	OnUse      func() error
	OnCancel   func()
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
