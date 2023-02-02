// rate package contains abstraction over the rate-limiting concept which is a strategy for limiting mostly network traffic.
// It allows putting a cap on how often someone can repeat an action within a certain timeframe.
// The package does not contain any certain implementation since multiple algorithms exist that are already implemented
// by the Golang community. Use this package to abstract away the specific implementations provided by third party libraries.
// It still would require to write an adapter for implementation to fulfill rate interfaces in your infrastracture layer.
package rate

// Limiter controls how frequently events are allowed to happen. The implementation decides which algorithm should be used
// that can fulfill the interface. The Token name is syntactic and does not restrict implementation to use
// token/leaky bucket algorithms. The interface serves as a simple API to rate limiting and any other algorithms
// like a fixed window can be used.
type Limiter interface {
	// Take returns a new Token. This should not consume the token and not consider the Token in availability calculation.
	// If Limit() is zero then this should return rate.FalseToken.
	Take() Token

	// Tokens should return the remaining token amount that can be consumed at now time with Take, so higher
	// Tokens value allow more events to happen without a delay. A zero value means none token can be consumed.
	Tokens() uint64

	// Limit should return the maximum amount of token that can be consumed within defined period with Take, so higher Limit
	// value allow more events to happen without a limit delay. A zero value means none token can be consumed and Take
	// should return rate.FalseToken.
	Limit() uint64
}

// BurstLimiter extends Limiter functionality to accept burst, so it adds an additional limit to allow multiple events
// to be consumed at once.
type BurstLimiter interface {
	Limiter

	// TakeN returns a new Token that allows to consume n tokens at once. This function does not consume the token
	// and should not consider the Token in availability calculation. If n is higher than Burst() it should return
	// rate.FalseToken.
	TakeN(n uint64) Token

	// Burst is the maximum number of tokens that can be consumed in a single call to TakeN, so higher Burst
	// value allow more events to happen at once. A zero value means none token can be consumed and TakeN
	// should always return rate.FalseToken.
	Burst() uint64
}

// ReservationLimiter extends Limiter functionality with reservations that are cancellable tokens so the caller
// can reserve a token which affects rate limiting calculation and decide later if he want to consume it or cancel
// and free up token allocation.
type ReservationLimiter interface {
	Limiter

	Reserve() CancellableToken
}

// BurstLimiter extends BurstLimiter and ReservationLimiter functionality to accept burst reservation, so it allow
// to reserve multiple tokens at once with single composed token.
type BurstReservationLimiter interface {
	BurstLimiter
	ReservationLimiter

	ReserveN(n uint64) CancellableToken
}
