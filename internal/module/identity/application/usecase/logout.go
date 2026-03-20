// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/command"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type Logout struct {
	sessions repository.SessionRepository
	mapper   mapper.AuthMapper
}

func NewLogout(
	session repository.SessionRepository,
	mapper mapper.AuthMapper,
) Logout {
	return Logout{session, mapper}
}

func (uc Logout) Execute(ctx context.Context, input dto.LogoutCommand) error {
	norm, err := uc.mapper.LogoutInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid logout data: %w", err)
	}
	if err := uc.sessions.Revoke(ctx, norm.IDSession); err != nil {
		return err
	}
	return nil
}
