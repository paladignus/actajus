package persistence

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/domainerrors"
)

// Mock usando função customizável
type MockAuthenticateFlex struct {
	SignInFunc func(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error)
	CallCount  int
}

func NewMockAuthenticateFlex() *MockAuthenticateFlex {
	return &MockAuthenticateFlex{
		SignInFunc: func(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
			return dto.AuthenticatedOutput{
				ID:        "1",
				FirstName: "Default",
				LastName:  "User",
			}, nil
		},
	}
}

func (m *MockAuthenticateFlex) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
	m.CallCount++
	return m.SignInFunc(ctx, cpf, password)
}

func TestAuthentication_CustomBehavior(t *testing.T) {
	mockAuth := NewMockAuthenticateFlex()
	mockAuth.SignInFunc = func(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
		if cpf == "special" {
			return dto.AuthenticatedOutput{ID: "999", FirstName: "Special", LastName: "User"}, nil
		}
		return dto.AuthenticatedOutput{}, domainerrors.ErrUserNotFound
	}

	// Teste com CPF especial
	result, err := mockAuth.SignIn(context.Background(), "special", "any")
	if err != nil {
		t.Error("Expected no error, got:", err)
	}
	if result.ID != "999" {
		t.Error("Expected ID 999, got:", result.ID)
	}
	// assert.NoError(t, err)
	// assert.Equal(t, 999, result.ID)

	// Teste com CPF normal
	result, err = mockAuth.SignIn(context.Background(), "normal", "any")
	if !errors.Is(err, domainerrors.ErrUserNotFound) {
		t.Error("Expected user not found error, got:", err)
	}
	// assert.ErrorIs(t, err, domainerrors.ErrUserNotFound)
}

// func TestAuthenticate(t *testing.T) {
// 	ctx := context.Background()
// 	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
// 	db, err := postgres.NewConnection(ctx, &config.DatabaseConfig{
// 		Host:     "localhost",
// 		Port:     5432,
// 		User:     "postgres",
// 		Password: "M4rc3l0",
// 		DBName:   "actajus_test",
// 		SSLMode:  "disable",
// 	}, logger)
// 	if err != nil {
// 		logger.Error("failed to connect to database", slog.Any("error", err))
// 	}
// 	sut := persistence.NewAuthenticate(db, logger)
// 	t.Run("should return user not found", func(t *testing.T) {
// 		_, err := sut.SignIn(ctx, "72775351115", "1234569")
// 		if err != domainerrors.ErrUserNotFound {
// 			t.Errorf("expected user not found error, got %v", err)
// 		}
// 	})
// 	t.Run("should authenticate user", func(t *testing.T) {
// 		user, err := sut.SignIn(ctx, "72775351115", "123456")
// 		if err != nil {
// 			t.Errorf("expected no error, got %v", err)
// 		}
// 		if user.ID == "" {
// 			t.Errorf("expected user id, got empty")
// 		}
// 		if user.FirstName != "Marcelo" {
// 			t.Errorf("expected first name Marcelo, got %s", user.FirstName)
// 		}
// 		if user.LastName != "Bento Pereira" {
// 			t.Errorf("expected last name Bento Pereira, got %s", user.LastName)
// 		}
// 	})
// 	db.Close()
// }
