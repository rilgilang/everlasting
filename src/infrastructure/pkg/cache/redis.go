package cache

import (
	"context"
	"encoding/json"
	"time"

	"everlasting/src/domain/sharedkernel/cache"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	values, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return cache.ErrCacheKeyNotFound
		}
		return err
	}
	err = json.Unmarshal([]byte(values), value)
	return err
}

func (c *RedisCache) Incr(ctx context.Context, key string) int {
	counts, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0
		}
		return 0
	}
	return int(counts)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, string(bytes), expiration).Err()
}

func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *RedisCache) Exist(ctx context.Context, key string) (int64, error) {
	return c.client.Exists(ctx, key).Result()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, cache.ErrCacheKeyNotFound
		}
		return 0, err
	}
	return ttl, nil
}
