// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type RecoverPassword interface {
	Execute(context.Context, dto.RecoverPasswordInput) error
}
