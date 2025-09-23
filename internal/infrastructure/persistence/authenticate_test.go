package persistence

import (
	"context"
	"log/slog"
	"os"
	"testing"

	errors "github.com/paladignus/actajus/internal/domain/error"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
)

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := postgres.NewConnection(ctx, &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "M4rc3l0",
		DBName:   "actajus_test",
		SSLMode:  "disable",
	}, logger)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
	}
	sut := NewAuthenticate(db, logger)
	t.Run("should return user not found", func(t *testing.T) {
		_, err := sut.SignIn(ctx, "72775351115", "1234569")
		if err != errors.ErrUserNotFound {
			t.Errorf("expected user not found error, got %v", err)
		}
	})
	t.Run("should authenticate user", func(t *testing.T) {
		user, err := sut.SignIn(ctx, "72775351115", "123456")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user.ID == "" {
			t.Errorf("expected user id, got empty")
		}
		if user.FirstName != "Marcelo" {
			t.Errorf("expected first name Marcelo, got %s", user.FirstName)
		}
		if user.LastName != "Bento Pereira" {
			t.Errorf("expected last name Bento Pereira, got %s", user.LastName)
		}
	})
	db.Close()
}
