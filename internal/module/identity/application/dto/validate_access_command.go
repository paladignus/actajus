// Package dto
package dto

type ValidateAccessCommand struct {
	AccessToken string `json:"access_token" validate:"required|min=10"`
}
