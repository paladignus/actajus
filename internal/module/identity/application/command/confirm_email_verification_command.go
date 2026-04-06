// Package command
package command

type ConfirmEmailVerificationCommand struct {
	IDVerification int64  `json:"id_verification" validate:"required|min=1"`
	Token          string `json:"token" validate:"required|min=10"`
}
