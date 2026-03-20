// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/command"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
)

type GrantPermissionToRole struct {
	roleUsers repository.RoleUserQueryRepository
	permRole  repository.PermissionRoleAdminRepository
	cache     service.RBACCacheInvalidator
	mapper    *mapper.RBACAdminMapper
}

func NewGrantPermissionToRole(
	roleUsers repository.RoleUserQueryRepository,
	permRole repository.PermissionRoleAdminRepository,
	cache service.RBACCacheInvalidator,
	mapper *mapper.RBACAdminMapper,
) GrantPermissionToRole {
	return GrantPermissionToRole{
		roleUsers,
		permRole,
		cache,
		mapper,
	}
}

func (uc GrantPermissionToRole) Execute(ctx context.Context, input dto.GrantPermissionToRoleCommand) error {
	norm, err := uc.mapper.GrantPermInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid grant permission data: %w", err)
	}
	if err := uc.permRole.GrantPermission(ctx, norm.IDRole, norm.IDPermission); err != nil {
		return err
	}
	idUsers, err := uc.roleUsers.ListUserIDsByRole(ctx, norm.IDRole)
	if err != nil {
		return err
	}
	for _, uid := range idUsers {
		uc.cache.InvalidateUser(ctx, uid)
	}
	return nil
}
