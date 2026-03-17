// Package dto
package dto

type UpdateEmailRequest struct {
	ID      int64  `json:"id_email" validate:"required"`
	Address string `json:"address" validate:"required|email"`
}
