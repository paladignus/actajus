// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
)

type SignIn struct {
	ExpectedOutput readmodel.SignInReadModel
	ExpectedError  error
}

func (m *SignIn) Execute(ctx context.Context, input command.SignInCommand) (readmodel.SignInReadModel, error) {
	return m.ExpectedOutput, m.ExpectedError
}
