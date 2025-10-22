// Package adapter
package adapter

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisCache(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewRedisCache(ctx, client)
	assert.NotNil(t, cache)
	assert.IsType(t, redisCache{}, cache)
}
