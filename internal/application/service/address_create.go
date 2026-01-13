// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type IAddressCreate interface {
	Execute(ctx context.Context, input dto.Address) error
}
