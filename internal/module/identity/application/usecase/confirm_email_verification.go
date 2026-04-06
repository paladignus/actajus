// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/command"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	identitydomain "github.com/paladignus/actajus/internal/module/identity/domain"
)

type ConfirmEmailVerification struct {
	user    repository.UserRepository
	verify  repository.EmailVerificationRepository
	refresh service.RefreshTokenService
	clock   service.Clock
	mapper  mapper.AuthMapper
}

func NewConfirmEmailVerification(
	user repository.UserRepository,
	verify repository.EmailVerificationRepository,
	refresh service.RefreshTokenService,
	clock service.Clock,
	mapper mapper.AuthMapper,
) ConfirmEmailVerification {
	return ConfirmEmailVerification{
		user:    user,
		verify:  verify,
		refresh: refresh,
		clock:   clock,
		mapper:  mapper,
	}
}

func (uc ConfirmEmailVerification) Execute(ctx context.Context, input command.ConfirmEmailVerificationCommand) (*readmodel.ConfirmEmailVerificationReadModel, error) {
	norm, err := uc.mapper.ConfirmEmailVerificationInputToNormalized(input)
	if err != nil {
		return nil, fmt.Errorf("invalid email verification data: %w", err)
	}
	now := uc.clock.Now()
	rec, err := uc.verify.GetByID(ctx, norm.IDVerification)
	if err != nil || rec == nil {
		return nil, identitydomain.ErrEmailVerificationNotFound
	}
	if rec.UsedAt != nil {
		return nil, identitydomain.ErrEmailVerificationUsed
	}
	if !now.Before(rec.ExpiresAt) {
		_ = uc.verify.MarkUsed(ctx, norm.IDVerification, now)
		return nil, identitydomain.ErrEmailVerificationExpired
	}
	if ok := uc.refresh.Compare(norm.Token, rec.Hash); !ok {
		return nil, identitydomain.ErrEmailVerificationInvalid
	}
	user, err := uc.user.FindByID(ctx, rec.IDUser)
	if err != nil {
		return nil, err
	}
	if err := user.MarkPrimaryEmailVerified(now); err != nil {
		return nil, err
	}
	if err := uc.user.MarkPrimaryEmailVerified(ctx, rec.IDEmail, now); err != nil {
		return nil, err
	}
	if err := uc.verify.MarkUsed(ctx, norm.IDVerification, now); err != nil {
		return nil, err
	}
	return &readmodel.ConfirmEmailVerificationReadModel{
		Message: "Email verificado com sucesso. Agora voce pode entrar no sistema.",
	}, nil
}
