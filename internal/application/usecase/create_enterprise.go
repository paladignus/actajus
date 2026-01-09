// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type Enterprise struct {
	enterprise repository.IEnterprise
}

func NewEnterprise(enterprise repository.IEnterprise) Enterprise {
	return Enterprise{enterprise}
}

func (e Enterprise) Execute(ctx context.Context, input dto.EnterpriseInput) error {
	enterprise, err := entity.NewEnterprise(input)
	if err != nil {
		return fmt.Errorf("use case create enterprise, invalid input: %w", err)
	}
	if err := e.enterprise.Create(ctx, enterprise); err != nil {
		return fmt.Errorf("database error while saving enterprise: %w", err)
	}
	return nil
}
