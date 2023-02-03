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

	// FalseToken should be returned from Take/TakeN methods to compare with IsFalseToken function.
	FalseToken CancellableToken = &falseToken{}
)

// Token is a controlled token by Limiter. It allows the caller to consume it or check if it can be consumed and
// the actual time that it'll be allowed to be consumed.
type Token interface {
	// Use consumes token and returns no error if succeeded. Returns ErrRateLimitExceeded when limit rate was exceeded.
	// If used with BurstLimiter it will not consume tokens below the limit before returning error. It will return ErrInvalidTokenValue
	// when Limiter Limit/Burst is lower than value provided in Take/TakeN.
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

// FalseToken is token returned when either tokens limit is 0 or reserved/taken token count is higher than burst/limit.
// Which means Limiter would never allow such token to be consumed.
type falseToken struct{}

// Allow always returns false.
func (*falseToken) Allow() bool {
	return false
}

// ResetsAt returns MaxTime value.
func (*falseToken) ResetsAt() time.Time {
	return MaxTime
}

// Use always returns ErrInvalidTokenValue.
func (*falseToken) Use() error {
	return ErrInvalidTokenValue
}

// Cancel is noop.
func (*falseToken) Cancel() {}

// IsFalseToken reports if token is always false meaning Limiter would never allow such token to be consumed.
func IsFalseToken(token Token) bool {
	return token == FalseToken
}
