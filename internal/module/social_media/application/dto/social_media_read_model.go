// Package dto
package dto

type SocialMediaReadModel struct {
	ID        uint   `json:"id"`
	IDCompany uint   `json:"id_company"`
	Platform  string `json:"platform"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
