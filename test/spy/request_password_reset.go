// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
)

type RequestPasswordReset struct {
	ExpectedError error
}

func (m *RequestPasswordReset) Execute(ctx context.Context, input command.RequestPasswordResetCommand) error {
	return m.ExpectedError
}
