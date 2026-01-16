// Package postgres
package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type UnitOfWork struct {
	db *pgxpool.Pool
}

func NewUnitOfWork(db *pgxpool.Pool) UnitOfWork {
	return UnitOfWork{db}
}

func (u UnitOfWork) Company() repository.UnitOfWorkCompany {
	return NewUnitOfWorkCompany(u.db)
}
