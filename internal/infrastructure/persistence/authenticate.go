// Package persistence
package persistence

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/domainerrors"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
)

type Authenticate struct {
	db *postgres.DB
}

func NewAuthenticate(db *postgres.DB) Authenticate {
	return Authenticate{db}
}

func (a Authenticate) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
	autenticated := dto.AuthenticatedOutput{}
	sql := `SELECT idpeople, first_name, last_name FROM people WHERE cpf = $1 AND password = crypt($2, password);`
	err := a.db.Pool.QueryRow(ctx, sql, cpf, password).
		Scan(&autenticated.UserID, &autenticated.FirstName, &autenticated.LastName)
	if errors.Is(err, pgx.ErrNoRows) {
		return autenticated, domainerrors.ErrUserNotFound
	}
	return autenticated, err
}
