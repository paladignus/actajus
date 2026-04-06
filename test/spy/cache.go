// Package spy
package spy

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type Cache struct {
	mock.Mock
}

func (m *Cache) Set(key string, value string, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *Cache) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *Cache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *Cache) Exists(key string) (int64, error) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *Cache) Close() error {
	args := m.Called()
	return args.Error(0)
}
