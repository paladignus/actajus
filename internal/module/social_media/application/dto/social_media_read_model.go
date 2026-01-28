// Package dto
package dto

type SocialMediaReadModel struct {
	ID        uint   `json:"id"`
	Platform  string `json:"platform"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
