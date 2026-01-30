// Package dto
package dto

type CreateSocialMediaRequest struct {
	Platform string `json:"platform" validate:"required"`
	URL      string `json:"url" validate:"required|url"`
}
