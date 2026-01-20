// Package spy
package spy

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

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
