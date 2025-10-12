// Package persistence
package persistence

import "github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"

type Persistence struct {
	db *postgres.DB
}

func NewPersistence(db *postgres.DB) Persistence {
	return Persistence{db}
}

func (p Persistence) Account() Account {
	return NewAccount(p.db)
}
