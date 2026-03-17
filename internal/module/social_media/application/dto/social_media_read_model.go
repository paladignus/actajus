// Package dto
package dto

import "time"

type SocialMediaReadModel struct {
	ID        int64     `json:"id"`
	IDCompany int64     `json:"id_company"`
	Platform  string    `json:"platform"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
