// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type SignIn interface {
	Execute(context.Context, dto.SignInInput) (dto.SignInOutput, error)
}
