package rate

type LimiterStore interface {
	Limit(key string) Limiter
}
