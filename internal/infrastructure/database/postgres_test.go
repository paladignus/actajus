// Package database
package database

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/test/infrastructure/adapter/spy"
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
