// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type ChangePassword struct {
	user              repository.UserRepository
	session           repository.SessionRepository
	hasher            service.PasswordHasher
	clock             service.Clock
	mapper            *mapper.AuthMapper
	RevokeAllSessions bool
}

func NewChangePassword(
	user repository.UserRepository,
	session repository.SessionRepository,
	hasher service.PasswordHasher,
	clock service.Clock,
	mapper *mapper.AuthMapper,
	revokeAll bool,
) ChangePassword {
	return ChangePassword{
		user, session, hasher, clock,
		mapper, revokeAll,
	}
}

func (uc ChangePassword) Execute(ctx context.Context, input dto.ChangePasswordCommand) error {
	norm, err := uc.mapper.ChangePasswordInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid change password data: %w", err)
	}
	uid := identity.IDUser(norm.IDUser)
	user, err := uc.user.FindByID(ctx, uid)
	if err != nil {
		return identity.ErrInvalidToken
	}
	if user.IsBlocked() {
		return identity.ErrUserBlocked
	}
	if err := uc.hasher.Compare(user.PasswordHash().Value(), norm.CurrentPassword); err != nil {
		return identity.ErrInvalidCredentials
	}
	newHash, err := uc.hasher.Hash(norm.NewPassword)
	if err != nil {
		return err
	}
	if err = uc.user.UpdatePasswordHash(ctx, uid, newHash); err != nil {
		return err
	}
	if uc.RevokeAllSessions {
		_ = uc.session.RevokeAllByUser(ctx, uid)
	}
	return nil
}
