package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache struct {
	keyPrefix string
	client    *redis.Client
}

func NewRedisCache(keyPrefix string, client *redis.Client) *RedisCache {
	return &RedisCache{
		keyPrefix: keyPrefix,
		client:    client,
	}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, c.withPrefix(key)).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist.
	}
	return val, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	return c.client.Set(ctx, c.withPrefix(key), value, time.Duration(expiration)*time.Second).Err()
}

func (c *RedisCache) FlushAll(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}

func (c *RedisCache) withPrefix(s string) string {
	return fmt.Sprintf("%s:%s", c.keyPrefix, s)
}
