// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type IEnterpriseGetAll interface {
	Execute(ctx context.Context) (output []dto.EnterpriseOutput, err error)
}
