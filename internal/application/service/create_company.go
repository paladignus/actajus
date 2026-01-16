// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type ICreateCompany interface {
	Execute(ctx context.Context, input dto.CompanyInputOutput) error
}
