// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/identity/application/command"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
)

type RBACAdminMapper struct {
	v *validation.Validator
}

func NewRBACAdminMapper(v *validation.Validator) *RBACAdminMapper {
	return &RBACAdminMapper{v: v}
}

type AssignRoleNormalized struct {
	IDUser     int64
	IDRole     int16
	AssignedBy int64
}

type RemoveRoleNormalized struct {
	IDUser int64
	IDRole int16
}

type GrantPermNormalized struct {
	IDRole       int16
	IDPermission int16
}

type RevokePermNormalized struct {
	IDRole       int16
	IDPermission int16
}

func (m *RBACAdminMapper) AssignRoleInputToNormalized(input command.AssignRoleToUserCommand) (AssignRoleNormalized, error) {
	vs := m.v.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return AssignRoleNormalized{}, err
	}
	return AssignRoleNormalized{
		IDUser:     input.IDUser,
		IDRole:     input.IDRole,
		AssignedBy: input.AssignedBy,
	}, nil
}

func (m *RBACAdminMapper) RemoveRoleInputToNormalized(input command.RemoveRoleFromUserCommand) (RemoveRoleNormalized, error) {
	vs := m.v.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return RemoveRoleNormalized{}, err
	}
	return RemoveRoleNormalized{
		IDUser: input.IDUser,
		IDRole: input.IDRole,
	}, nil
}

func (m *RBACAdminMapper) GrantPermInputToNormalized(input command.GrantPermissionToRoleCommand) (GrantPermNormalized, error) {
	vs := m.v.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return GrantPermNormalized{}, err
	}
	return GrantPermNormalized{
		IDRole:       input.IDRole,
		IDPermission: input.IDPermission,
	}, nil
}

func (m *RBACAdminMapper) RevokePermInputToNormalized(input command.RevokePermissionFromRoleCommand) (RevokePermNormalized, error) {
	vs := m.v.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return RevokePermNormalized{}, err
	}
	return RevokePermNormalized{
		IDRole:       input.IDRole,
		IDPermission: input.IDPermission,
	}, nil
}
