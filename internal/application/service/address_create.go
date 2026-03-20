// Package service
package service

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
)

type IAddressCreate interface {
	Execute(ctx context.Context, input command.CreateAddressCommand) error
}
