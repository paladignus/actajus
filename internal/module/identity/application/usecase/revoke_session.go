// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/command"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type RevokeSession struct {
	sessions repository.SessionRepository
	mapper   mapper.AuthMapper
}

func NewRevokeSession(
	session repository.SessionRepository,
	mapper mapper.AuthMapper,
) RevokeSession {
	return RevokeSession{session, mapper}
}

func (uc RevokeSession) Execute(ctx context.Context, input dto.RevokeSessionCommand) error {
	norm, err := uc.mapper.RevokeInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid revoke data: %w", err)
	}
	if err := uc.sessions.Revoke(ctx, norm.IDSession); err != nil {
		return err
	}
	return nil
}
