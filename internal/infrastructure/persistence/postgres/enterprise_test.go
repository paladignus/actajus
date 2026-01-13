// Package postgres
package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestEnterprise(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()
	// db := database.DB{Pool: mock}
	repo := NewEnterprise(mock)
	input := dto.Enterprise{RegisteredBy: 1, Name: "Actajus", TradeName: "Actajus Trade", CNPJ: "10.123.456/0001-00"}
	enterprise, err := entity.NewEnterprise(input)
	assert.NoError(t, err)

	t.Run("should be successfully inseted into the database", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO companies`).
			WithArgs(enterprise.RegisteredBy, enterprise.Name, enterprise.TradeName, enterprise.CNPJ).
			WillReturnRows(pgxmock.NewRows([]string{"idcompanies"}).AddRow("1"))
		id, err := repo.Create(ctx, enterprise)
		assert.NoError(t, err)
		assert.Equal(t, "1", id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when database fails to insert", func(t *testing.T) {
		mockErr := fmt.Errorf("database error")
		mock.ExpectQuery(`INSERT INTO companies`).
			WithArgs(enterprise.RegisteredBy, enterprise.Name, enterprise.TradeName, enterprise.CNPJ).
			WillReturnError(mockErr)
		id, err := repo.Create(ctx, enterprise)
		assert.Error(t, err)
		assert.Empty(t, id)
		assert.ErrorContains(t, err, "database error while saving enterprise")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
