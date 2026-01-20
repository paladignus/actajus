// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type RenewPassword interface {
	Execute(ctx context.Context, input dto.RenewPasswordInput) error
}
