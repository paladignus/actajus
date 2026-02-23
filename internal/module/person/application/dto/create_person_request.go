// Package dto
package dto

type CreatePersonRequest struct {
	Name     string `json:"name" validate:"required|min=2"`
	Birthday string `json:"birthday" validate:"required"`
	Gender   uint   `json:"gender" validate:"required|oneof=1,2"`
	Mother   uint   `json:"mother"`
	Father   uint   `json:"father"`
}
