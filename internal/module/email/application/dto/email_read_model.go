// Package dto
package dto

type EmailReadModel struct {
	ID        uint   `json:"id"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
