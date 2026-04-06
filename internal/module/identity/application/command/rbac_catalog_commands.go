// Package command
package command

type CreateRoleCommand struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateRoleCommand struct {
	ID          int16  `json:"id" validate:"required|min=1"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type DeleteRoleCommand struct {
	ID int16 `json:"id" validate:"required|min=1"`
}

type CreatePermissionCommand struct {
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdatePermissionCommand struct {
	ID          int16  `json:"id" validate:"required|min=1"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type DeletePermissionCommand struct {
	ID int16 `json:"id" validate:"required|min=1"`
}
