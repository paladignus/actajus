// Package persistence
package persistence

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/application/dto"
	errors "github.com/paladignus/actajus/internal/domain/error"
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
