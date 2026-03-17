// Package dto
package dto

type UpdatePhoneRequest struct {
	ID         int64  `json:"id_phone" validate:"required|numeric"`
	Number     string `json:"number" validate:"required|min=13|max=14"`
	Kind       string `json:"kind" validate:"required"`
	Department string `json:"department,omitzero"`
}
