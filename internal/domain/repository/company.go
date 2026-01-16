// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
)

type ICompany interface {
	Create(ctx context.Context, enterprise entity.Company) (id uint, err error)
	Update(ctx context.Context, enterprise entity.Company) error
	GetAll(ctx context.Context) ([]dto.CompanyInputOutput, error)
	Delete(ctx context.Context, idCompany uint) error
}
