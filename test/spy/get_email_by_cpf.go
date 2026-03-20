// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
)

type GetEmailByCPF struct {
	ExpectedOutput readmodel.GetEmailByCPFReadModel
	ExpectedError  error
}

func (m *GetEmailByCPF) Execute(ctx context.Context, input command.GetEmailByCPFCommand) (readmodel.GetEmailByCPFReadModel, error) {
	return m.ExpectedOutput, m.ExpectedError
}
