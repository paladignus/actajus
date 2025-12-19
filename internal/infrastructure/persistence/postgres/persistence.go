// Package postgres
package postgres

import (
	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type Persistence struct {
	db *database.DB
}

func NewPersistence(db *database.DB) Persistence {
	return Persistence{db}
}

func (p Persistence) User() User {
	return NewUser(p.db)
}

func (p Persistence) PasswordResetToken() PasswordResetToken {
	return NewPasswordResetToken(p.db)
}
