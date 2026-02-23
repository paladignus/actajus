// Package usecase
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/model"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
)

type PasswordResetConfig struct {
	ResetTTL time.Duration
}

type RequestPasswordReset struct {
	user           repository.UserRepository
	reset          repository.PasswordResetRepository
	refresh        service.RefreshTokenService // reutiliza gerador/hash (mesmo do refresh!)
	clock          service.Clock
	cfg            PasswordResetConfig
	mapper         mapper.AuthMapper
	RevokePrevious bool
}

func NewRequestPasswordReset(
	user repository.UserRepository,
	reset repository.PasswordResetRepository,
	refresh service.RefreshTokenService,
	clock service.Clock,
	cfg PasswordResetConfig,
	mapper mapper.AuthMapper,
	revokePrevious bool,
) RequestPasswordReset {
	return RequestPasswordReset{
		user, reset, refresh, clock,
		cfg, mapper, revokePrevious,
	}
}

func (uc RequestPasswordReset) Execute(ctx context.Context, input dto.RequestPasswordResetCommand) (*dto.RequestPasswordResetReadModel, error) {
	norm, err := uc.mapper.RequestPasswordResetInputToNormalized(input)
	if err != nil {
		return nil, fmt.Errorf("invalid password reset request data: %w", err)
	}
	user, err := uc.user.FindByEmail(ctx, norm.Email)
	if err != nil || user == nil {
		return &dto.RequestPasswordResetReadModel{
			IDReset:    0,
			ResetToken: "",
			ExpiresAt:  time.Time{},
		}, nil
	}
	now := uc.clock.Now()
	if uc.RevokePrevious {
		_ = uc.reset.RevokeAllByUser(ctx, user.ID(), now)
	}
	token, hash, err := uc.refresh.Generate()
	if err != nil {
		return nil, err
	}
	expiresAt := now.Add(uc.cfg.ResetTTL)
	in := model.PasswordResetTokenCreate{
		IDUser:    user.ID(),
		Hash:      hash,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}
	rid, err := uc.reset.Create(ctx, &in)
	if err != nil {
		return nil, err
	}
	return &dto.RequestPasswordResetReadModel{
		IDReset:    rid.Value(),
		ResetToken: token,
		ExpiresAt:  expiresAt,
	}, nil
}
