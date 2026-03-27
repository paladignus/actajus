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

type RevokePermissionFromRole struct {
	roleUsers repository.RoleUserQueryRepository
	permRole  repository.PermissionRoleAdminRepository
	cache     service.RBACCacheInvalidator
	mapper    *mapper.RBACAdminMapper
}

func NewRevokePermissionFromRole(
	roleUsers repository.RoleUserQueryRepository,
	permRole repository.PermissionRoleAdminRepository,
	cache service.RBACCacheInvalidator,
	mapper *mapper.RBACAdminMapper,
) RevokePermissionFromRole {
	return RevokePermissionFromRole{
		roleUsers,
		permRole,
		cache,
		mapper,
	}
}

func (uc RevokePermissionFromRole) Execute(ctx context.Context, input command.RevokePermissionFromRoleCommand) error {
	norm, err := uc.mapper.RevokePermInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid revoke permission data: %w", err)
	}
	if err := uc.permRole.RevokePermission(ctx, norm.IDRole, norm.IDPermission); err != nil {
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
