package rate

import (
	"errors"
	"time"
)

var (
	// ErrRateLimitExceeded is the error returned by Token.Use when the rate limit exceeds the limit.
	ErrRateLimitExceeded = errors.New("maximum rate limit exceeded")

	// ErrInvalidTokenValue is the error returned by Token.Use when Token is always falsy.
	ErrInvalidTokenValue = errors.New("token value exceeds limit")

	// MinTime is minimum time value.
	MinTime = time.Unix(-2208988800, 0)

	// MaxTime is maximum time value.
	MaxTime = MinTime.Add(1<<63 - 1)
)

// Token is a controlled token by Limiter. It allows the caller to consume it or check if it can be consumed and
// the actual time that it'll be allowed to be consumed.
type Token interface {
	// Use consumes token and returns no error if succeeded. Returns ErrRateLimitExceeded when limit rate was exceeded.
	// If used with BurstLimiter it will not consume tokens below the limit before returning error. It will return ErrInvalidTokenValue
	// when Limit in BurstLimiter is lower than value provided in TakeN.
	Use() error

	// Allow reports whether an token can be consumed.
	Allow() bool

	// ResetsAt returns a time when token will be available to be consumed. At returned time Allow() should report true.
	ResetsAt() time.Time
}

// CancellableToken is a controlled token by ReservationLimiter. It extends the Token functionality with Cancel function
// that allows to free up the reserved tokens.
type CancellableToken interface {
	Token

	// Cancel frees up tokens to ReservationLimiter.
	Cancel()
}
