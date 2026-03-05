// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
)

type AssignRoleToUser struct {
	repo   identityrepo.RoleUserAdminRepository
	cache  service.RBACCacheInvalidator
	index  service.RBACRoleUsersIndex
	mapper *mapper.RBACAdminMapper
}

func NewAssignRoleToUser(
	repo identityrepo.RoleUserAdminRepository,
	cache service.RBACCacheInvalidator,
	index service.RBACRoleUsersIndex,
	mapper *mapper.RBACAdminMapper,
) AssignRoleToUser {
	return AssignRoleToUser{
		repo:   repo,
		cache:  cache,
		index:  index,
		mapper: mapper,
	}
}

func (uc AssignRoleToUser) Execute(ctx context.Context, input dto.AssignRoleToUserCommand) error {
	norm, err := uc.mapper.AssignRoleInputToNormalized(input)
	if err != nil {
		return fmt.Errorf("invalid assign role data: %w", err)
	}
	if err := uc.repo.AssignRole(ctx, norm.IDUser, norm.IDRole, norm.AssignedBy); err != nil {
		return err
	}
	uc.index.AddUserToRole(ctx, norm.IDRole, norm.IDUser)
	uc.cache.InvalidateUser(ctx, norm.IDUser)
	return nil
}
