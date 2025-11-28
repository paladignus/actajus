// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type RecoverPassword struct {
	ExpectedError error
}

func (m *RecoverPassword) Execute(ctx context.Context, input dto.RecoverPasswordInput) error {
	return m.ExpectedError
}
