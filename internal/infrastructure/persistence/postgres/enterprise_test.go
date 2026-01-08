// Package postgres
package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/infrastructure/database"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

type Enterprise struct {
	db *database.DB
}

func NewEnterprise(db *database.DB) *Enterprise {
	return &Enterprise{db: db}
}

func (e Enterprise) Create(ctx context.Context, enterprise entity.Enterprise) error {
	sql := `INSERT INTO enterprises (name, trade_name, cnpj) VALUES ($1, $2, $3);`
	_, err := e.db.Pool.Exec(ctx, sql, enterprise.Name, enterprise.TradeName, enterprise.CNPJ)
	if err != nil {
		return fmt.Errorf("database error while saving enterprise: %w", err)
	}
	return nil
}

func TestEnterprise(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()
	db := &database.DB{Pool: mock}
	repo := NewEnterprise(db)
	enterprise, err := entity.NewEnterprise("Actajus", "Actajus Trade", "10.123.456/0001-00")
	assert.NoError(t, err)
	// enterprise := entity.Enterprise{
	// 	Name:      "actajus",
	// 	TradeName: "actajus-trade",
	// 	CNPJ:      "10.123.456/0001-00",
	// }

	t.Run("should be successfully inseted into the database", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO enterprises`).
			WithArgs(enterprise.Name, enterprise.TradeName, enterprise.CNPJ).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		err := repo.Create(ctx, enterprise)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when database fails to insert", func(t *testing.T) {
		mockErr := fmt.Errorf("database error")
		mock.ExpectExec(`INSERT INTO enterprises`).
			WithArgs(enterprise.Name, enterprise.TradeName, enterprise.CNPJ).
			WillReturnError(mockErr)
		err := repo.Create(ctx, enterprise)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database error while saving enterprise")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
