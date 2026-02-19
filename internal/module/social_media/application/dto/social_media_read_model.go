// Package dto
package dto

import "time"

type SocialMediaReadModel struct {
	ID        uint      `json:"id"`
	IDCompany uint      `json:"id_company"`
	Platform  string    `json:"platform"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
