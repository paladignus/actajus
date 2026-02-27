// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
)

type RemoveRoleFromUser struct {
	repo   repository.RoleUserAdminRepository
	cache  service.RBACCacheInvalidator
	mapper *mapper.RBACAdminMapper
}

func NewRemoveRoleFromUser(
	repo repository.RoleUserAdminRepository,
	cache service.RBACCacheInvalidator,
	mapper *mapper.RBACAdminMapper,
) RemoveRoleFromUser {
	return RemoveRoleFromUser{repo, cache, mapper}
}

func (uc RemoveRoleFromUser) Execute(ctx context.Context, input dto.RemoveRoleFromUserCommand) error {
	norm, err := uc.mapper.RemoveRoleInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid remove role data: %w", err)
	}
	if err := uc.repo.RemoveRole(ctx, norm.IDUser, norm.IDRole); err != nil {
		return err
	}

	uc.cache.InvalidateUser(ctx, norm.IDUser)
	return nil
}
