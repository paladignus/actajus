// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type IUpdateCompany interface {
	Execute(ctx context.Context, input dto.CompanyInputOutput) error
}
