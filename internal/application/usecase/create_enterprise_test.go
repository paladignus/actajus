// Package usecase
package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestCreateEnterprise(t *testing.T) {
	ctx := context.Background()
	repo := spy.NewEnterprise()
	sut := NewEnterprise(repo)
	input := dto.EnterpriseInput{
		RegisteredBy: 1, Name: "Actajus", TradeName: "Actajus Trade", CNPJ: "10.123.456/0001-00",
	}

	t.Run("should create a new enterprise", func(t *testing.T) {
		err := sut.Execute(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("should return error when database fails to insert", func(t *testing.T) {
		repoErr := fmt.Errorf("database error")
		repo.CreateError = repoErr
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database error while saving enterprise")
		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("should return error when input is invalid", func(t *testing.T) {
		input.RegisteredBy = 0
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "use case create enterprise, invalid input")
	})
}
