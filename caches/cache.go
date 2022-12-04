// Package caches provides abstraction for both caching strategies in memory and distributed.
package caches

import (
	"context"
	"time"
)

// Cache is an interface defining cache contract with individual ttl configuration on set.
type Cache interface {
	// Set should store any value with ttl support.
	Set(key string, value interface{}, ttl time.Duration) error

	// Get should assign cached value to v which must be a pointer.
	Get(key string, v interface{}) error

	// Delete should erase/invalidate cache with corresponding key.
	Delete(key string) error
}

// DistributedCache is an interface defining cache which can be distributed on server like redis.
// The different between cache and distributed in interface is accepting context to respect
// requests cancellation.
type DistributedCache interface {
	// Set should store any value with ttl support.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Get should assign cached value to v which must be a pointer.
	Get(ctx context.Context, key string, v interface{}) error

	// Delete should erase/invalidate cache with corresponding key.
	Delete(ctx context.Context, key string) error
}
