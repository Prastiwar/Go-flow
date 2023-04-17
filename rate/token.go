package rate

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrRateLimitExceeded is the error returned by Token.Use when the rate limit exceeds the limit.
	ErrRateLimitExceeded = errors.New("maximum rate limit exceeded")

	// ErrInvalidTokenValue is the error returned by Token.Use when Token is always falsy.
	ErrInvalidTokenValue = errors.New("token value exceeds limit")

	// ErrTokenCancelled is the error returned by Token.Use when CancellableToken.Cancel was used.
	ErrTokenCancelled = errors.New("token value exceeds limit")

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

	// Allow reports whether an token can be consumed at this moment.
	Allow() bool

	// ResetsAt returns a time when token will be available to be consumed. At returned time Allow() should report true.
	ResetsAt() time.Time

	// Context returns the token's context. To change the context, use WithContext.
	// The returned context should be always non-nil; it defaults to the background context.
	// Since the token is controlled by Limiter, the context should be the value passed in Limiter.Take.
	// The token's context is immutable. To change it create a new Token using Limiter.Take.
	Context() context.Context
}

// CancellableToken is a controlled token by ReservationLimiter. It extends the Token functionality with Cancel function
// that allows to free up the reserved tokens. In opposite to Token this should not be reusable and should be used
// one-time only for either consumption or cancellation. If token was canceled, Use function should return ErrTokenCancelled.
type CancellableToken interface {
	Token

	// Cancel frees up tokens to ReservationLimiter without actually consuming the token. Cancellation does not
	// allow to consume the token and Use function will return ErrTokenCancelled error.
	Cancel()
}
