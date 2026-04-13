package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrCacheKeyNotFound error = errors.New("cache key not found")
)

type (
	CacheRepository interface {
		Get(ctx context.Context, key string, value interface{}) error
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Delete(ctx context.Context, keys ...string) error
		Exist(ctx context.Context, key string) (int64, error)
		Close() error
		TTL(ctx context.Context, key string) (time.Duration, error)
		Incr(ctx context.Context, key string) int
	}
)
