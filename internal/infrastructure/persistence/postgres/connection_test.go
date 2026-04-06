// Package postgres
package postgres

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConnection_Success(t *testing.T) {
	ctx := context.Background()
	mockLogger := spy.Logger{}
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "M4rc3l0",
		DBName:   "actajus",
		SSLMode:  "disable",
	}
	t.Run("should return a new database connection", func(t *testing.T) {
		mockLogger.On("Info", ctx, "database connection sucessfully").Once()
		db, err := NewConnection(ctx, cfg, &mockLogger)
		require.NoError(t, err)
		require.NotNil(t, db)
		// require.NotNil(t, db.Pool)
		mockLogger.On("Info", ctx, "database connection closed").Once()
		// db.Close(ctx, &mockLogger)
		db.Close()
		mockLogger.AssertExpectations(t)
	})

	t.Run("should return error on failure from malformed DSN", func(t *testing.T) {
		cfg = &config.DatabaseConfig{
			Host:     "host with spaces",
			Port:     "5432",
			User:     "user with @",
			Password: "pass with ://",
			DBName:   "db name",
			SSLMode:  "disable",
		}
		db, err := NewConnection(ctx, cfg, &mockLogger)
		assert.Error(t, err)
		assert.True(t,
			containsAny(err.Error(), "unable to parse database config", "unable to create connection pool"),
			"Error should be about parsing or connection: %v", err)
		assert.Nil(t, db)
	})

	t.Run("should return error on failure from ping", func(t *testing.T) {
		cfg = &config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "invalid",
			DBName:   "actajus",
			SSLMode:  "disable",
		}
		_, err := NewConnection(ctx, cfg, &mockLogger)
		assert.Error(t, err)
		assert.True(t,
			containsAny(err.Error(), "unable to connect to database", "unable to connect to database"),
			"Error should be about connection: %v", err)
	})

	// t.Run("should closed connection", func(t *testing.T) {
	// 	mockPool := spy.MockPgxPool{}
	// 	mockLogger := spy.Logger{}
	// 	db := &DB{Pool: &mockPool}
	// 	mockPool.On("Close").Once()
	// 	mockLogger.On("Info", ctx, "database connection closed").Once()
	// 	db.Close(ctx, &mockLogger)
	// 	mockPool.AssertExpectations(t)
	// 	mockLogger.AssertExpectations(t)
	// })
}

func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if contains(s, substr) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}

// func TestNewConnection_PingError(t *testing.T) {
// 	originalParseConfig := pgxpoolParseConfig
// 	originalNewWithConfig := pgxpoolNewWithConfig
// 	defer func() {
// 		pgxpoolParseConfig = originalParseConfig
// 		pgxpoolNewWithConfig = originalNewWithConfig
// 	}()
// 	ctx := context.Background()
// 	mockLogger := spy.Logger{}
// 	validConfig := &config.DatabaseConfig{
// 		Host:     "localhost",
// 		Port:     "5432",
// 		User:     "user",
// 		Password: "pass",
// 		DBName:   "db",
// 		SSLMode:  "disable",
// 	}
// 	mockPool := &MockPgxPool{}
// 	mockPool.On("Ping", ctx).Return(errors.New("forced ping error"))
// 	// ParseConfig funciona
// 	pgxpoolParseConfig = func(connString string) (*pgxpool.Config, error) {
// 		return &pgxpool.Config{}, nil
// 	}
// 	// NewWithConfig retorna pool mockado
// 	pgxpoolNewWithConfig = func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
// 		// Precisamos converter nosso mock para *pgxpool.Pool
// 		// Como não podemos fazer isso diretamente, vamos criar um pool real que vai falhar
// 		// usando uma configuração inválida
// 		return pgxpool.New(ctx, "postgres://invalid:invalid@invalid:9999/invalid")
// 	}
// 	db, err := NewConnection(ctx, validConfig, &mockLogger)
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "unable to connect to database")
// 	assert.Nil(t, db)
// }
