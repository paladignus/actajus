// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type LogoutAll struct {
	session repository.SessionRepository
	mapper  mapper.AuthMapper
}

func NewLogoutAll(
	session repository.SessionRepository,
	mapper mapper.AuthMapper,
) LogoutAll {
	return LogoutAll{session, mapper}
}

func (uc LogoutAll) Execute(ctx context.Context, input dto.LogoutAllCommand) error {
	norm, err := uc.mapper.LogoutAllInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid logout-all data: %w", err)
	}
	if err := uc.session.RevokeAllByUser(ctx, norm.IDUser); err != nil {
		return err
	}

	return nil
}
