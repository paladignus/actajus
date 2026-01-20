// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type SignIn struct {
	ExpectedOutput dto.SignInOutput
	ExpectedError  error
}

func (m *SignIn) Execute(ctx context.Context, input dto.SignInInput) (dto.SignInOutput, error) {
	return m.ExpectedOutput, m.ExpectedError
}
