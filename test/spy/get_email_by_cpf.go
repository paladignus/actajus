// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type GetEmailByCPF struct {
	ExpectedOutput dto.GetEmailByCPFOutput
	ExpectedError  error
}

func (m *GetEmailByCPF) Execute(ctx context.Context, input dto.GetEmailByCPFInput) (dto.GetEmailByCPFOutput, error) {
	return m.ExpectedOutput, m.ExpectedError
}
