// Package dto
package dto

type UpdateEmailRequest struct {
	IDEmail int64  `json:"id_email" validate:"required"`
	Address string `json:"address" validate:"required|email"`
}
