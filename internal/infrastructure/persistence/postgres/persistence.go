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

func (p Persistence) Authentication() Authentication {
	return NewAuthentication(p.db)
}

func (p Persistence) Token() Token {
	return NewToken(p.db)
}
