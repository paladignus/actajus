// Package database
package database

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/test/infrastructure/adapter/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock para substituir pgxpool.ParseConfig durante os testes
var pgxpoolParseConfig = pgxpool.ParseConfig

// Mock para substituir pgxpool.NewWithConfig durante os testes
var pgxpoolNewWithConfig = pgxpool.NewWithConfig

// Restante dos mocks anteriores (MockPgxPool, MockLogger) permanecem iguais...
type MockPgxPool struct {
	mock.Mock
}

func (m *MockPgxPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.Rows), argsMock.Error(1)
}

func (m *MockPgxPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.Row)
}

func (m *MockPgxPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgconn.CommandTag), argsMock.Error(1)
}

func (m *MockPgxPool) Close() {
	m.Called()
}

func TestNewConnection_Success(t *testing.T) {
	ctx := context.Background()
	mockLogger := spy.SpyLogger{}
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "M4rc3l0",
		DBName:   "actajus",
		SSLMode:  "disable",
	}
	mockLogger.On("Info", mock.Anything, "database connection sucessfully").Once()
	db, err := NewConnection(ctx, cfg, &mockLogger)
	require.NoError(t, err)
	require.NotNil(t, db)
	require.NotNil(t, db.Pool)
	mockLogger.On("Info", mock.Anything, "database connection closed").Once()
	db.Close(ctx, &mockLogger)
	mockLogger.AssertExpectations(t)
}

func TestNewConnection_MalformedDSN(t *testing.T) {
	ctx := context.Background()
	mockLogger := spy.SpyLogger{}
	// Configuração que gera DSN malformada
	config := &config.DatabaseConfig{
		Host:     "host with spaces", // Host com espaços - gera DSN inválida
		Port:     "5432",
		User:     "user with @",   // User com @ - caracter especial
		Password: "pass with ://", // Password com :// - quebra a URL
		DBName:   "db name",       // DB com espaço
		SSLMode:  "disable",
	}
	db, err := NewConnection(ctx, config, &mockLogger)
	assert.Error(t, err)
	// Pode ser parse error ou connection error dependendo de como o pgx lida
	assert.True(t,
		containsAny(err.Error(), "unable to parse database config", "unable to create connection pool"),
		"Error should be about parsing or connection: %v", err)
	assert.Nil(t, db)
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

func TestNewConnection_PingError(t *testing.T) {
	originalParseConfig := pgxpoolParseConfig
	originalNewWithConfig := pgxpoolNewWithConfig
	defer func() {
		pgxpoolParseConfig = originalParseConfig
		pgxpoolNewWithConfig = originalNewWithConfig
	}()
	ctx := context.Background()
	mockLogger := spy.SpyLogger{}
	validConfig := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
	}
	mockPool := &MockPgxPool{}
	mockPool.On("Ping", mock.Anything).Return(errors.New("forced ping error"))
	// ParseConfig funciona
	pgxpoolParseConfig = func(connString string) (*pgxpool.Config, error) {
		return &pgxpool.Config{}, nil
	}
	// NewWithConfig retorna pool mockado
	pgxpoolNewWithConfig = func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		// Precisamos converter nosso mock para *pgxpool.Pool
		// Como não podemos fazer isso diretamente, vamos criar um pool real que vai falhar
		// usando uma configuração inválida
		return pgxpool.New(ctx, "postgres://invalid:invalid@invalid:9999/invalid")
	}
	db, err := NewConnection(ctx, validConfig, &mockLogger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to connect to database")
	assert.Nil(t, db)
}

func TestCloseDB(t *testing.T) {
	ctx := context.Background()
	mockPool := MockPgxPool{}
	mockLogger := spy.SpyLogger{}
	db := &DB{Pool: &mockPool}
	mockPool.On("Close").Once()
	mockLogger.On("Info", mock.Anything, "database connection closed").Once()
	db.Close(ctx, &mockLogger)
	mockPool.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
