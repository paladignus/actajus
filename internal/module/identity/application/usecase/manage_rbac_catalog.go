// Package usecase
package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/paladignus/actajus/internal/module/identity/application/command"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
)

type CreateRole struct {
	repo      repository.RoleAdminRepository
	validator *validation.Validator
}

func NewCreateRole(repo repository.RoleAdminRepository, validator *validation.Validator) CreateRole {
	return CreateRole{repo: repo, validator: validator}
}

func (uc CreateRole) Execute(ctx context.Context, input command.CreateRoleCommand) error {
	if err := sharedAdapter.ViolationsToDomainError(uc.validator.ValidateStruct(input)); err != nil {
		return err
	}
	return uc.repo.Create(ctx, strings.TrimSpace(input.Name), strings.TrimSpace(input.Description))
}

type UpdateRole struct {
	repo      repository.RoleAdminRepository
	validator *validation.Validator
}

func NewUpdateRole(repo repository.RoleAdminRepository, validator *validation.Validator) UpdateRole {
	return UpdateRole{repo: repo, validator: validator}
}

func (uc UpdateRole) Execute(ctx context.Context, input command.UpdateRoleCommand) error {
	if err := sharedAdapter.ViolationsToDomainError(uc.validator.ValidateStruct(input)); err != nil {
		return err
	}
	return uc.repo.Update(ctx, input.ID, strings.TrimSpace(input.Name), strings.TrimSpace(input.Description))
}

type DeleteRole struct {
	repo repository.RoleAdminRepository
}

func NewDeleteRole(repo repository.RoleAdminRepository) DeleteRole {
	return DeleteRole{repo: repo}
}

func (uc DeleteRole) Execute(ctx context.Context, input command.DeleteRoleCommand) error {
	if input.ID <= 0 {
		return fmt.Errorf("invalid delete role data")
	}
	return uc.repo.Delete(ctx, input.ID)
}

type CreatePermission struct {
	repo      repository.PermissionAdminRepository
	validator *validation.Validator
}

func NewCreatePermission(repo repository.PermissionAdminRepository, validator *validation.Validator) CreatePermission {
	return CreatePermission{repo: repo, validator: validator}
}

func (uc CreatePermission) Execute(ctx context.Context, input command.CreatePermissionCommand) error {
	if err := sharedAdapter.ViolationsToDomainError(uc.validator.ValidateStruct(input)); err != nil {
		return err
	}
	return uc.repo.Create(
		ctx,
		strings.TrimSpace(input.Resource),
		strings.TrimSpace(input.Action),
		strings.TrimSpace(input.Description),
	)
}

type UpdatePermission struct {
	repo      repository.PermissionAdminRepository
	validator *validation.Validator
}

func NewUpdatePermission(repo repository.PermissionAdminRepository, validator *validation.Validator) UpdatePermission {
	return UpdatePermission{repo: repo, validator: validator}
}

func (uc UpdatePermission) Execute(ctx context.Context, input command.UpdatePermissionCommand) error {
	if err := sharedAdapter.ViolationsToDomainError(uc.validator.ValidateStruct(input)); err != nil {
		return err
	}
	return uc.repo.Update(
		ctx,
		input.ID,
		strings.TrimSpace(input.Resource),
		strings.TrimSpace(input.Action),
		strings.TrimSpace(input.Description),
	)
}

type DeletePermission struct {
	repo repository.PermissionAdminRepository
}

func NewDeletePermission(repo repository.PermissionAdminRepository) DeletePermission {
	return DeletePermission{repo: repo}
}

func (uc DeletePermission) Execute(ctx context.Context, input command.DeletePermissionCommand) error {
	if input.ID <= 0 {
		return fmt.Errorf("invalid delete permission data")
	}
	return uc.repo.Delete(ctx, input.ID)
}
