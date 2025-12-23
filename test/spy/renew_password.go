// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type RenewPassword struct {
	ExpectedError error
}

func (r *RenewPassword) Execute(ctx context.Context, input dto.RenewPasswordInput) error {
	return r.ExpectedError
}
