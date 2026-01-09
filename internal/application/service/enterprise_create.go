// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type IEnterpriseCreate interface {
	Execute(ctx context.Context, input dto.EnterpriseInput) error
}
