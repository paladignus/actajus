// Package adapter
package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestRedisCache_Get(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available:", err)
	}
	cache := NewRedisCache(ctx, client)
	defer client.Close()
	tests := []struct {
		name      string
		key       string
		setValue  string
		wantValue string
		setup     bool
		wantErr   error
	}{
		{
			name:      "should be successful when getting",
			key:       "test:get:1",
			setValue:  "test-value",
			wantValue: "test-value",
			setup:     true,
			wantErr:   nil,
		},
		{
			name:      "should return empty for a non-existent value with error",
			key:       "test:get:nonexistent",
			setValue:  "",
			wantValue: "",
			setup:     false,
			wantErr:   redis.Nil,
		},
		{
			name:      "should return empty for a non-existent value without error",
			key:       "test:get:2",
			setValue:  "",
			wantValue: "",
			setup:     true,
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup {
				err := cache.Set(tt.key, tt.setValue, 10*time.Second)
				require.NoError(t, err)
				defer cache.Delete(tt.key)
			}
			value, err := cache.Get(tt.key)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, value)
			}
		})
	}
}
