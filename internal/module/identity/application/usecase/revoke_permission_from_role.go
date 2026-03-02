// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	"github.com/paladignus/actajus/internal/module/identity/domain"
)

type RevokePermissionFromRole struct {
	roleUsers repository.RoleUserAdminRepository
	permRole  repository.PermissionRoleAdminRepository
	cache     service.RBACCacheInvalidator
	index     service.RBACRoleUsersIndex
	mapper    *mapper.RBACAdminMapper
}

func NewRevokePermissionFromRole(
	roleUsers repository.RoleUserAdminRepository,
	permRole repository.PermissionRoleAdminRepository,
	cache service.RBACCacheInvalidator,
	index service.RBACRoleUsersIndex,
	mapper *mapper.RBACAdminMapper,
) RevokePermissionFromRole {
	return RevokePermissionFromRole{
		roleUsers,
		permRole,
		cache,
		index,
		mapper,
	}
}

func (uc RevokePermissionFromRole) Execute(ctx context.Context, input dto.RevokePermissionFromRoleCommand) error {
	norm, err := uc.mapper.RevokePermInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid revoke permission data: %w", err)
	}
	if err := uc.permRole.RevokePermission(ctx, norm.IDRole, norm.IDPermission); err != nil {
		return err
	}
	idUsers, err := resolveRoleUsersWithFallback(ctx, norm.IDRole, uc.index, uc.roleUsers)
	if err != nil {
		return err
	}
	for _, uid := range idUsers {
		uc.cache.InvalidateUser(ctx, domain.IDUser(uid))
	}
	return nil
}
