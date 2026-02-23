// Package dto
package dto

type RequestPasswordResetCommand struct {
	Email string `json:"email" validate:"required|email"`
}
