// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
)

type RequestPasswordReset interface {
	Execute(ctx context.Context, input command.RequestPasswordResetCommand) error
}
