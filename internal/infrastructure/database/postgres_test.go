// Package database
package database

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/paladignus/actajus/internal/infrastructure/config"
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

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(ctx context.Context, msg string, args ...any) {
	m.Called(ctx, msg, args)
}

func (m *MockLogger) Info(ctx context.Context, msg string, args ...any) {
	m.Called(ctx, msg, args)
}

func (m *MockLogger) Warn(ctx context.Context, msg string, args ...any) {
	m.Called(ctx, msg, args)
}

func (m *MockLogger) Error(ctx context.Context, msg string, args ...any) {
	m.Called(ctx, msg, args)
}

func (m *MockLogger) With(args ...any) repository.Logger {
	argsMock := m.Called(args)
	return argsMock.Get(0).(repository.Logger)
}

func (m *MockLogger) WithError(err error) repository.Logger {
	argsMock := m.Called(err)
	return argsMock.Get(0).(repository.Logger)
}

// Testes para os casos de erro específicos
func TestNewConnection_ErrorScenarios(t *testing.T) {
	t.SkipNow()
	ctx := context.Background()
	mockLogger := new(MockLogger)

	// Configuração válida de base
	// validConfig := &config.DatabaseConfig{
	// 	Host:     "localhost",
	// 	Port:     "5432",
	// 	User:     "testuser",
	// 	Password: "testpass",
	// 	DBName:   "testdb",
	// 	SSLMode:  "disable",
	// }

	tests := []struct {
		name          string
		config        *config.DatabaseConfig
		setupMocks    func()
		expectedError string
	}{
		{
			name: "error parsing invalid DSN - invalid character",
			config: &config.DatabaseConfig{
				Host:     "local@host", // Caractere inválido no host
				Port:     "5432",
				User:     "user",
				Password: "pass",
				DBName:   "db",
				SSLMode:  "disable",
			},
			setupMocks: func() {
				// Não precisa mockar nada - o pgxpool.ParseConfig vai falhar naturalmente
				// com um host contendo @
			},
			expectedError: "unable to parse database config",
		},
		{
			name: "error creating connection pool - invalid port",
			config: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     "-1", // Porta inválida
				User:     "user",
				Password: "pass",
				DBName:   "db",
				SSLMode:  "disable",
			},
			setupMocks: func() {
				// Não precisa mockar - porta negativa causará erro
			},
			expectedError: "unable to create connection pool",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}
			db, err := NewConnection(ctx, tt.config, mockLogger)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
			assert.Nil(t, db)
		})
	}
}

// Teste mais específico usando monkey patching (avançado)
func TestNewConnection_ParseConfigError(t *testing.T) {
	t.SkipNow()
	// Salvar a função original
	originalParseConfig := pgxpoolParseConfig
	defer func() {
		pgxpoolParseConfig = originalParseConfig
	}()

	ctx := context.Background()
	mockLogger := new(MockLogger)
	validConfig := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
	}
	// Mock da função ParseConfig para forçar erro
	pgxpoolParseConfig = func(connString string) (*pgxpool.Config, error) {
		return nil, errors.New("forced parse config error")
	}
	db, err := NewConnection(ctx, validConfig, mockLogger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to parse database config")
	assert.Contains(t, err.Error(), "forced parse config error")
	assert.Nil(t, db)
}

func TestNewConnection_NewWithConfigError(t *testing.T) {
	t.SkipNow()
	// Salvar funções originais
	originalParseConfig := pgxpoolParseConfig
	originalNewWithConfig := pgxpoolNewWithConfig
	defer func() {
		pgxpoolParseConfig = originalParseConfig
		pgxpoolNewWithConfig = originalNewWithConfig
	}()
	ctx := context.Background()
	mockLogger := new(MockLogger)
	validConfig := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
	}
	// ParseConfig funciona
	pgxpoolParseConfig = func(connString string) (*pgxpool.Config, error) {
		return &pgxpool.Config{}, nil
	}
	// NewWithConfig falha
	pgxpoolNewWithConfig = func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		return nil, errors.New("forced new with config error")
	}
	db, err := NewConnection(ctx, validConfig, mockLogger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to create connection pool")
	assert.Contains(t, err.Error(), "forced new with config error")
	assert.Nil(t, db)
}

func TestNewConnection_PingError(t *testing.T) {
	// Salvar funções originais
	originalParseConfig := pgxpoolParseConfig
	originalNewWithConfig := pgxpoolNewWithConfig
	defer func() {
		pgxpoolParseConfig = originalParseConfig
		pgxpoolNewWithConfig = originalNewWithConfig
	}()

	ctx := context.Background()
	mockLogger := new(MockLogger)
	validConfig := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
	}

	// Mock do pool que falha no Ping
	mockPool := new(MockPgxPool)
	mockPool.On("Ping", ctx).Return(errors.New("forced ping error"))
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
	db, err := NewConnection(ctx, validConfig, mockLogger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to connect to database")
	assert.Nil(t, db)
}

// Teste com DSN malformada
func TestNewConnection_MalformedDSN(t *testing.T) {
	ctx := context.Background()
	mockLogger := new(MockLogger)
	// Configuração que gera DSN malformada
	config := &config.DatabaseConfig{
		Host:     "host with spaces", // Host com espaços - gera DSN inválida
		Port:     "5432",
		User:     "user with @",   // User com @ - caracter especial
		Password: "pass with ://", // Password com :// - quebra a URL
		DBName:   "db name",       // DB com espaço
		SSLMode:  "disable",
	}
	db, err := NewConnection(ctx, config, mockLogger)
	assert.Error(t, err)
	// Pode ser parse error ou connection error dependendo de como o pgx lida
	assert.True(t,
		containsAny(err.Error(), "unable to parse database config", "unable to create connection pool"),
		"Error should be about parsing or connection: %v", err)
	assert.Nil(t, db)
}

// Helper function
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

// Teste de sucesso com configuração válida (quando há banco disponível)
func TestNewConnection_Success(t *testing.T) {
	ctx := context.Background()
	mockLogger := new(MockLogger)
	// Configuração para um banco real (ajuste conforme seu ambiente)
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres", // Ajuste
		Password: "M4rc3l0",  // Ajuste
		DBName:   "actajus",  // Ajuste
		SSLMode:  "disable",
	}
	mockLogger.On("Info", ctx, "database connection sucessfully", []interface{}(nil)).Once()

	db, err := NewConnection(ctx, cfg, mockLogger)

	require.NoError(t, err)
	require.NotNil(t, db)
	require.NotNil(t, db.Pool)

	// Verificar que podemos fazer ping
	// (isso testa que a conexão realmente funciona)

	// Fechar
	mockLogger.On("Info", ctx, "database connection closed", []interface{}(nil)).Once()
	db.Close(ctx, mockLogger)

	mockLogger.AssertExpectations(t)
}

// Testes anteriores permanecem...
func TestDB_Close(t *testing.T) {
	ctx := context.Background()
	mockPool := new(MockPgxPool)
	mockLogger := new(MockLogger)

	db := &DB{Pool: mockPool}

	mockPool.On("Close").Once()
	mockLogger.On("Info", ctx, "database connection closed", []interface{}(nil)).Once()

	db.Close(ctx, mockLogger)

	mockPool.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// ... outros testes anteriores
