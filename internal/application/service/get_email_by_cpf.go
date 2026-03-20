// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
)

type GetEmailByCPF interface {
	Execute(dto context.Context, input command.GetEmailByCPFCommand) (output readmodel.GetEmailByCPFReadModel, err error)
}
