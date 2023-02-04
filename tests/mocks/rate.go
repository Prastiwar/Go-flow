package mocks

import (
	"time"

	"github.com/Prastiwar/Go-flow/rate"
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
	return m.OnLimit(key)
}

type LimiterMock struct {
	OnLimit  func() uint64
	OnTokens func() uint64
	OnTake   func() rate.Token
}

func (m LimiterMock) Limit() uint64 {
	return m.OnLimit()
}

func (m LimiterMock) Take() rate.Token {
	return m.OnTake()
}

func (m LimiterMock) Tokens() uint64 {
	return m.OnTokens()
}

type BurstLimiterMock struct {
	rate.Limiter

	OnBurst func() uint64
	OnTakeN func(n uint64) rate.Token
}

func (m BurstLimiterMock) TakeN(n uint64) rate.Token {
	return m.OnTakeN(n)
}

func (m BurstLimiterMock) Burst() uint64 {
	return m.OnBurst()
}

type TokenMock struct {
	OnAllow    func() bool
	OnResetsAt func() time.Time
	OnUse      func() error
	OnCancel   func()
}

func (m TokenMock) Allow() bool {
	return m.OnAllow()
}

func (m TokenMock) ResetsAt() time.Time {
	return m.OnResetsAt()
}

func (m TokenMock) Use() error {
	return m.OnUse()
}

func (m TokenMock) Cancel() {
	m.OnCancel()
}
