// Package caches provides abstraction for both caching strategies in memory and distributed.
package caches

import (
	"context"
	"time"
)

// Cache is an interface defining evicting cache contract with individual ttl configuration.
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Evict(evict func(key string, data []byte)) error
}

// DistributedCache is an interface defining cache which can be distributed on server like redis.
// The different between cache and distributed in interface is accepting context to respect
// requests cancellation.
type DistributedCache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}
