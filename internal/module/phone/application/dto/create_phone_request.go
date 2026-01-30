// Package dto
package dto

type CreatePhoneRequest struct {
	Number     string `json:"number" validate:"required|min=13|max=14"`
	Kind       string `json:"kind" validate:"required"`
	Department string `json:"department,omitzero"`
}
