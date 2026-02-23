// Package dto
package dto

type PersonReadModel struct {
	ID        uint   `json:"id"`
	Gender    uint32 `json:"gender"`
	Name      string `json:"name"`
	Birthday  string `json:"birthday"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
