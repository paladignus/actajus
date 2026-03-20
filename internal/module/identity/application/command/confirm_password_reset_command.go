// Package dto
package dto

type ConfirmPasswordResetCommand struct {
	IDReset     int64  `json:"id_reset" validate:"required|min=1"`
	ResetToken  string `json:"reset_token" validate:"required|min=10"`
	NewPassword string `json:"new_password" validate:"required|min=8"`
}
