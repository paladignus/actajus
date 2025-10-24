// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type Authenticate interface {
	Execute(context.Context, dto.AuthenticateInput) (dto.AuthenticateOutput, error)
}
