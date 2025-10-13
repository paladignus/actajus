// Package persistence
package persistence

import (
	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type Persistence struct {
	db *database.DB
}

func NewPersistence(db *database.DB) Persistence {
	return Persistence{db}
}

func (p Persistence) Account() Account {
	return NewAccount(p.db)
}
