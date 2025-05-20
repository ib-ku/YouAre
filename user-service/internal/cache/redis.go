package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisCache{client: rdb}
}

func (r *RedisCache) Set(key string, value string, ttl time.Duration) error {
	return r.client.Set(Ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(Ctx, key).Result()
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(Ctx, key).Err()
}
