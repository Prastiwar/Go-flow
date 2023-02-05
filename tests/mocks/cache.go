package mocks

import (
	"context"
	"time"

	"github.com/Prastiwar/Go-flow/caches"
)

var (
	_ caches.Cache            = CacheMock{}
	_ caches.DistributedCache = DistributedCacheMock{}
)

type CacheMock struct {
	OnDelete func(key string) error
	OnGet    func(key string, v interface{}) error
	OnSet    func(key string, value interface{}, ttl time.Duration) error
}

func (m CacheMock) Delete(key string) error {
	return m.OnDelete(key)
}

func (m CacheMock) Get(key string, v interface{}) error {
	return m.OnGet(key, v)
}

func (m CacheMock) Set(key string, value interface{}, ttl time.Duration) error {
	return m.OnSet(key, value, ttl)
}

type DistributedCacheMock struct {
	OnDelete func(ctx context.Context, key string) error
	OnGet    func(ctx context.Context, key string, v interface{}) error
	OnSet    func(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

func (m DistributedCacheMock) Delete(ctx context.Context, key string) error {
	return m.OnDelete(ctx, key)
}

func (m DistributedCacheMock) Get(ctx context.Context, key string, v interface{}) error {
	return m.OnGet(ctx, key, v)
}

func (m DistributedCacheMock) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return m.OnSet(ctx, key, value, ttl)
}
