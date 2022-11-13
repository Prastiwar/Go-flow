package cashes

type Cache interface {
	Set(ctx context.Context, key interface{}, value interface{}, expireTime time.Duration) error
	Get(ctx context.Context, key interface{}) ([]byte, error)
	Delete(ctx context.Context key interface{}) error
}