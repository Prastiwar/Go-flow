package rate

// LimiterStore is store for different Limiters indicated with key string.
type LimiterStore interface {
	// Limit returns Limiter that was persisted before in store with corresponding key string.
	// If there was no entry found it should return default Limiter and persist it with specified key.
	Limit(key string) Limiter
}