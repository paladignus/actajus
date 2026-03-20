// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/command"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type ConfirmPasswordReset struct {
	user              repository.UserRepository
	session           repository.SessionRepository
	reset             repository.PasswordResetRepository
	hasher            service.PasswordHasher
	refresh           service.RefreshTokenService
	clock             service.Clock
	mapper            mapper.AuthMapper
	RevokeAllSessions bool
}

func NewConfirmPasswordReset(
	user repository.UserRepository,
	session repository.SessionRepository,
	reset repository.PasswordResetRepository,
	hasher service.PasswordHasher,
	refresh service.RefreshTokenService,
	clock service.Clock,
	mapper mapper.AuthMapper,
	revokeAll bool,
) ConfirmPasswordReset {
	return ConfirmPasswordReset{
		user, session, reset,
		hasher, refresh, clock, mapper,
		revokeAll,
	}
}

func (uc ConfirmPasswordReset) Execute(ctx context.Context, input dto.ConfirmPasswordResetCommand) error {
	norm, err := uc.mapper.ConfirmPasswordResetInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid password reset confirm data: %w", err)
	}
	now := uc.clock.Now()
	rec, err := uc.reset.GetByID(ctx, norm.IDReset)
	if err != nil || rec == nil {
		return identity.ErrResetTokenNotFound
	}
	if rec.UsedAt != nil {
		return identity.ErrResetTokenUsed
	}
	if !now.Before(rec.ExpiresAt) {
		_ = uc.reset.MarkUsed(ctx, norm.IDReset, now)
		return identity.ErrResetTokenExpired
	}
	if ok := uc.refresh.Compare(norm.ResetToken, rec.Hash); !ok {
		return identity.ErrResetTokenInvalid
	}

	// Busca o usuário para usar o comportamento de domínio
	user, err := uc.user.FindByID(ctx, rec.IDUser)
	if err != nil {
		return err
	}

	// Delega a lógica de negócio para a entidade User
	// A entidade valida e gera o hash da nova senha
	if err := user.ChangePassword(norm.NewPassword, uc.hasher); err != nil {
		return err
	}

	// Persiste o hash atualizado
	if err := uc.user.UpdatePasswordHash(ctx, rec.IDUser, user.PasswordHash().Value()); err != nil {
		return err
	}
	if err := uc.reset.MarkUsed(ctx, norm.IDReset, now); err != nil {
		return err
	}
	if uc.RevokeAllSessions {
		_ = uc.session.RevokeAllByUser(ctx, rec.IDUser)
	}
	return nil
}
