// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
)

type IUpdateCompany interface {
	Execute(ctx context.Context, input command.UpdateCompanyCommand) error
}
