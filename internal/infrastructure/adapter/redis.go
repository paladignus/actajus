// Package adapter
package adapter

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisCache(ctx context.Context, client *redis.Client) repository.Cache {
	return redisCache{ctx, client}
}

func (r redisCache) Set(key, value string, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r redisCache) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r redisCache) Exists(key string) (int64, error) {
	return r.client.Exists(r.ctx, key).Result()
}

func (r redisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
