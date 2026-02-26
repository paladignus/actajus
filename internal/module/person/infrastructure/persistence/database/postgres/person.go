// Package postgres
package postgres

import (
	"context"

	"github.com/paladignus/actajus/internal/module/person/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type Person struct {
	db postgres.PgxPool
}

func NewPerson(db postgres.PgxPool) Person {
	return Person{db}
}

func (p Person) Create(ctx context.Context, person *domain.Person) error {
	query := `
		INSERT INTO people 
			(first_name, last_name, id_gender, birthday, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) returning idpeople`
	var id uint
	if err := p.db.QueryRow(ctx, query,
		person.FirstName(),
		person.LastName(),
		person.IDGender(),
		person.Birthday().Value(),
		person.CreatedAt(),
		person.UpdatedAt(),
	).Scan(&id); err != nil {
		return err
	}
	return person.SetID(id)
}
