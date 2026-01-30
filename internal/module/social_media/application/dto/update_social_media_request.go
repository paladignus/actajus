// Package dto
package dto

type UpdateSocialMediaRequest struct {
	IDSocialMedia uint   `json:"id_social_media" validate:"required|numeric"`
	Platform      string `json:"platform" validate:"required"`
	URL           string `json:"url" validate:"required|url"`
}
