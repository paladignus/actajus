// Package dto
package dto

type CreateSocialMediaRequest struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
}
