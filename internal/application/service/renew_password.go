// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
)

type RenewPassword interface {
	Execute(ctx context.Context, input command.RenewPasswordCommand) error
}
