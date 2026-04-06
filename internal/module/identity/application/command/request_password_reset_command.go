// Package command
package command

type RequestPasswordResetCommand struct {
	Email string `json:"email" validate:"required|email"`
}
