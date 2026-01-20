// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type GetEmailByCPF interface {
	Execute(dto context.Context, input dto.GetEmailByCPFInput) (output dto.GetEmailByCPFOutput, err error)
}
