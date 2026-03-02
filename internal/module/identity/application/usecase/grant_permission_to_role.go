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

type GrantPermissionToRole struct {
	roleUsers repository.RoleUserAdminRepository
	permRole  repository.PermissionRoleAdminRepository
	cache     service.RBACCacheInvalidator
	index     service.RBACRoleUsersIndex
	mapper    *mapper.RBACAdminMapper
}

func NewGrantPermissionToRole(
	roleUsers repository.RoleUserAdminRepository,
	permRole repository.PermissionRoleAdminRepository,
	cache service.RBACCacheInvalidator,
	index service.RBACRoleUsersIndex,
	mapper *mapper.RBACAdminMapper,
) GrantPermissionToRole {
	return GrantPermissionToRole{
		roleUsers,
		permRole,
		cache,
		index,
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
	idUsers, err := resolveRoleUsersWithFallback(ctx, norm.IDRole, uc.index, uc.roleUsers)
	if err != nil {
		return err
	}
	for _, uid := range idUsers {
		uc.cache.InvalidateUser(ctx, domain.IDUser(uid))
	}
	return nil
}
