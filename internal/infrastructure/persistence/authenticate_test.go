package persistence

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/application/dto"
	errors "github.com/paladignus/actajus/internal/domain/error"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
)

type Authenticate struct {
	db  *postgres.DB
	log *slog.Logger
}

func NewAuthenticate(db *postgres.DB, logger *slog.Logger) Authenticate {
	return Authenticate{db, logger}
}

func (a Authenticate) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
	autenticated := dto.AuthenticatedOutput{}
	sql := `SELECT idpeople, first_name, last_name FROM people WHERE cpf = $1 AND password = crypt($2, password);`
	err := a.db.Pool.QueryRow(ctx, sql, cpf, password).
		Scan(&autenticated.ID, &autenticated.FirstName, &autenticated.LastName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return autenticated, errors.ErrUserNotFound
		}
		a.log.Error("failed to autenticate", slog.Any("error", err))
	}
	return autenticated, err
}

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
	db.Close()
}
