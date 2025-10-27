// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type GetEmailByCPF interface {
	Execute(context.Context, dto.GetEmailByCPFInput) (dto.GetEmailByCPFOutput, error)
}
