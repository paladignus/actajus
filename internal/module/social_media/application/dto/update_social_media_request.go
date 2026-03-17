// Package dto
package dto

type UpdateSocialMediaRequest struct {
	IDSocialMedia int64  `json:"id_social_media" validate:"required|numeric"`
	Platform      string `json:"platform" validate:"required"`
	URL           string `json:"url" validate:"required|url"`
}
