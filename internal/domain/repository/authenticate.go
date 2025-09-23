// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type Authenticate interface {
	SignIn(ctx context.Context, cpf, password string) (dto.AuthenticatedOutput, error)
}
