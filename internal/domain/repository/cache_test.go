package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type CacheSpy struct {
	data map[string]string
}

func (c *CacheSpy) Set(key string, value string, duration time.Duration) error {
	if c.data == nil {
		c.data = make(map[string]string)
	}
	c.data[key] = value
	return nil
}

func (c *CacheSpy) Get(key string) (string, error) {
	if c.data == nil {
		return "", nil
	}
	value, exists := c.data[key]
	if !exists {
		return "", nil
	}
	return value, nil
}

func (c *CacheSpy) Exists(key string) (int64, error) {
	if c.data == nil {
		return 0, nil
	}
	_, exists := c.data[key]
	if exists {
		return 1, nil
	}
	return 0, nil
}

func (c *CacheSpy) Delete(key string) error {
	if c.data != nil {
		delete(c.data, key)
	}
	return nil
}

func TestCacheInterface(t *testing.T) {
	sut := &CacheSpy{}
	err := sut.Set("key1", "value1", time.Minute)
	assert.NoError(t, err)
	value, err := sut.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value)
	count, err := sut.Exists("key1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	err = sut.Delete("key1")
	assert.NoError(t, err)
	count, err = sut.Exists("key1")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

