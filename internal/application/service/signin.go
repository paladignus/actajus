// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type SignIn interface {
	Execute(ctx context.Context, input dto.SignInInput) (output dto.SignInOutput, err error)
}
