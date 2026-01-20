// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type RequestPasswordReset struct {
	ExpectedError error
}

func (m *RequestPasswordReset) Execute(ctx context.Context, input dto.RequestPasswordResetInput) error {
	return m.ExpectedError
}
