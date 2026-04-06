// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
)

type SignIn interface {
	Execute(ctx context.Context, input command.SignInCommand) (output readmodel.SignInReadModel, err error)
}
