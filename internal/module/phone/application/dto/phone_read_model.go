// Package dto
package dto

type PhoneReadModel struct {
	ID         string `json:"id"`
	Number     string `json:"number"`
	Kind       string `json:"kind"`
	Department string `json:"department"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
