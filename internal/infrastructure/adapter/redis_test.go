// Package adapter
package adapter

import (
	"context"
	"testing"
	"time"

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

func TestRedisCache_Set(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use separate DB for tests
	})
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available:", err)
	}
	cache := NewRedisCache(ctx, client)
	defer client.Close()
	tests := []struct {
		name    string
		key     string
		value   string
		ttl     time.Duration
		wantErr bool
	}{
		{
			name:    "should be successful when storing with expiration",
			key:     "test:set:1",
			value:   "test-value",
			ttl:     5 * time.Minute,
			wantErr: false,
		},
		{
			name:    "should be successful when storing without expiration",
			key:     "test:set:2",
			value:   "test-value",
			ttl:     0,
			wantErr: false,
		},
		{
			name:    "should be successful when storing empty value",
			key:     "test:set:3",
			value:   "",
			ttl:     5 * time.Minute,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer cache.Delete(tt.key)
			err := cache.Set(tt.key, tt.value, tt.ttl)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
