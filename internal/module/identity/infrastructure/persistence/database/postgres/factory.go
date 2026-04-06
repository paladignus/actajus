// Package postgres
package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type Factory struct {
	exec postgres.Executor
}

func NewFactory(exec postgres.Executor) *Factory {
	return &Factory{exec: exec}
}

func (f *Factory) WithTx(tx uow.Tx) repository.Factory {
	pgxTx, ok := tx.(pgx.Tx)
	if !ok {
		panic(fmt.Sprintf("invalid tx type: %T", tx))
	}
	return &Factory{exec: pgxTx}
}

func (f *Factory) PasswordReset() repository.PasswordResetRepository {
	return NewPasswordReset(f.exec)
}

func (f *Factory) EmailVerification() repository.EmailVerificationRepository {
	return NewEmailVerification(f.exec)
}

func (f *Factory) User() repository.UserRepository {
	return NewUser(f.exec)
}
