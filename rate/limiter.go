// rate package contains abstraction over the rate-limiting concept which is a strategy for limiting mostly network traffic.
// It allows putting a cap on how often someone can repeat an action within a certain timeframe.
// The package does not contain any certain implementation since multiple algorithms exist that are already implemented
// by the Golang community. Use this package to abstract away the specific implementations provided by third party libraries.
// It still would require to write an adapter for implementation to fulfill rate interfaces in your infrastracture layer.
package rate

import "context"

// Limiter controls how frequently events are allowed to happen. The implementation decides which algorithm should be used
// that can fulfill the interface. The Token name is syntactic and does not restrict implementation to use
// token/leaky bucket algorithms. The interface serves as a simple API to rate limiting and any other algorithms
// like a fixed window can be used.
type Limiter interface {
	// Take returns a new Token. This should not consume the token and not consider the Token in availability calculation.
	// Token should be reusable and time resistant and not usable only once it was got. The context will be passed down to
	// token which will use it if required for calculations.
	Take(ctx context.Context) Token

	// Tokens should return the remaining token amount that can be consumed at now time with Take, so higher
	// Tokens value allow more events to happen without a delay. A zero value means none token can be consumed now.
	Tokens(ctx context.Context) uint64

	// Limit should return the maximum amount of token that can be consumed within defined period with Take, so higher Limit
	// value allow more events to happen without a limit delay.
	Limit() uint64
}

// BurstLimiter extends Limiter functionality to accept burst, so it adds an additional limit to allow multiple events
// to be consumed at once.
type BurstLimiter interface {
	Limiter

	// TakeN returns a new Token that allows to consume n tokens at once. This should not consume the token and not consider
	// the Token in availability calculation. Token should be reusable and time resistant and not usable only once it was got.
	// The context will be passed down to token which will use it if required for calculations.
	TakeN(ctx context.Context, n uint64) Token

	// Burst is the maximum number of tokens that can be consumed in a single call to TakeN, so higher Burst
	// value allow more events to happen at once. This is not the maximum available value to be consumed. Use Tokens
	// if you need such value.
	Burst() uint64
}

// ReservationLimiter extends Limiter functionality with reservations that are cancellable tokens so the caller
// can reserve a token which affects rate limiting calculation and decide later if he want to consume it or cancel
// and free up token allocation.
type ReservationLimiter interface {
	Limiter

	// Reserve returns a new CancellableToken that allows to consume token. This should not consume the token
	// but would potentially consider in availability calculation. CancellableToken should not be reusable and used for one-time only.
	// The context will be passed down to token which will use it if required for calculations.
	Reserve(ctx context.Context) CancellableToken
}

// BurstLimiter extends BurstLimiter and ReservationLimiter functionality to accept burst reservation, so it allow
// to reserve multiple tokens at once with single composed token.
type BurstReservationLimiter interface {
	BurstLimiter
	ReservationLimiter

	// ReserveN returns a new CancellableToken that allows to consume n tokens at once. This should not consume the token
	// but would potentially consider in availability calculation. CancellableToken should not be reusable and used for one-time only.
	// The context will be passed down to token which will use it if required for calculations.
	ReserveN(ctx context.Context, n uint64) CancellableToken
}
