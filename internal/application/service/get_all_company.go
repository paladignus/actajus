// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type IGetAllCompany interface {
	Execute(ctx context.Context) (output []dto.CompanyInputOutput, err error)
}
