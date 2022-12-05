// Package caches provides an abstraction for both caching strategies in memory and distributed.
// Due to the complexity of implementing the correct caching strategy, you should make an adapter
// of a mature third-party solution in your infrastructure layer to fulfill the contract.
// The package contains both the contract and exported errors that are expected to
// be returned by implementation from specified methods.
package caches

import (
	"context"
	"errors"
	"time"
)

var (
	NoExpiration = time.Duration(0)

	ErrNotFound   = errors.New("cache value was not found for specified key")
	ErrNotPointer = errors.New("interface value must be a pointer")
)

// Cache is an interface defining cache contract with individual ttl configuration on set.
type Cache interface {
	// Set should store any value with ttl support. If it already exists - no error should be returned.
	// It should be replaced instead. In case value is replaced - it's up to implementation if eviction is supported.
	// Expected behaviour for ttl set to NoExpiration is that value should not expire and user must use Delete to remove it.
	Set(key string, value interface{}, ttl time.Duration) error

	// Get should assign cached value to v which must be a pointer - otherwise it should return an error.
	// If value was not set before, implementation should return ErrNotFound. In case v is nil or not a pointer
	// ErrNotPointer should be returned.
	Get(key string, v interface{}) error

	// Delete should erase/invalidate cache with corresponding key.
	// If value was not set before, implementation should return ErrNotFound.
	Delete(key string) error
}

// DistributedCache is an interface defining cache which can be distributed on server like redis.
// The difference between cache and distributed in interface is accepting context to respect
// requests cancellation.
type DistributedCache interface {
	// Set should store any value with ttl support. If it already exists - no error should be returned.
	// It should be replaced instead. In case value is replaced - it's up to implementation if eviction is supported.
	// Expected behaviour for ttl set to NoExpiration is that value should not expire and user must use Delete to remove it.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Get should assign cached value to v which must be a pointer - otherwise it should return an error.
	// If value was not set before, implementation should return ErrNotFound. In case v is nil or not a pointer
	// ErrNotPointer should be returned.
	Get(ctx context.Context, key string, v interface{}) error

	// Delete should erase/invalidate cache with corresponding key.
	// If value was not set before, implementation should return ErrNotFound.
	Delete(ctx context.Context, key string) error
}
