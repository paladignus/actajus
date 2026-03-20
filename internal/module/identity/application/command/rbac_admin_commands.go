// Package dto
package dto

type AssignRoleToUserCommand struct {
	IDUser     int64 `json:"id_user" validate:"required|min=1"`
	IDRole     int16 `json:"id_role" validate:"required|min=1"`
	AssignedBy int64 `json:"assigned_by" validate:"required|min=1"`
}

type RemoveRoleFromUserCommand struct {
	IDUser int64 `json:"id_user" validate:"required|min=1"`
	IDRole int16 `json:"id_role" validate:"required|min=1"`
}

type GrantPermissionToRoleCommand struct {
	IDRole       int16 `json:"id_role" validate:"required|min=1"`
	IDPermission int16 `json:"id_permission" validate:"required|min=1"`
}

type RevokePermissionFromRoleCommand struct {
	IDRole       int16 `json:"id_role" validate:"required|min=1"`
	IDPermission int16 `json:"id_permission" validate:"required|min=1"`
}
