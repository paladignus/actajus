// Package dto
package dto

type ChangePasswordCommand struct {
	IDUser          int64  `json:"id_user" validate:"required|min=1"`
	CurrentPassword string `json:"current_password" validate:"required|min=8"`
	NewPassword     string `json:"new_password" validate:"required|min=8"`
}
