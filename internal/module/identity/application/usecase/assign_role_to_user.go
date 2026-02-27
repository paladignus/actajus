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

type AssignRoleToUser struct {
	repo   repository.RoleUserAdminRepository
	cache  service.RBACCacheInvalidator
	mapper *mapper.RBACAdminMapper
}

func NewAssignRoleToUser(
	repo repository.RoleUserAdminRepository,
	cache service.RBACCacheInvalidator,
	mapper *mapper.RBACAdminMapper,
) AssignRoleToUser {
	return AssignRoleToUser{repo, cache, mapper}
}

func (uc AssignRoleToUser) Execute(ctx context.Context, input dto.AssignRoleToUserCommand) error {
	norm, err := uc.mapper.AssignRoleInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid assign role data: %w", err)
	}
	if err := uc.repo.AssignRole(ctx, norm.IDUser, norm.IDRole, norm.AssignedBy); err != nil {
		return err
	}
	uc.cache.InvalidateUser(ctx, norm.IDUser)
	return nil
}
