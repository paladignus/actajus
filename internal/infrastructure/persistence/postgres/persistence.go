// Package postgres
package postgres

import (
	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type Persistence struct {
	db database.PgxPool
}

func NewPersistence(db database.PgxPool) Persistence {
	return Persistence{db}
}

func (p Persistence) User() User {
	return NewUser(p.db)
}

func (p Persistence) PasswordResetToken() PasswordResetToken {
	return NewPasswordResetToken(p.db)
}

func (p Persistence) Enterprise() Enterprise {
	return *NewEnterprise(p.db)
}
