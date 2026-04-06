// Package command
package command

type RequestEmailVerificationCommand struct {
	Email string `json:"email" validate:"required|email"`
}
