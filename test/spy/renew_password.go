// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
)

type RenewPassword struct {
	ExpectedError error
}

func (r *RenewPassword) Execute(ctx context.Context, input command.RenewPasswordCommand) error {
	return r.ExpectedError
}
