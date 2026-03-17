// Package dto
package dto

import "time"

type PhoneReadModel struct {
	ID         int64     `json:"id"`
	Number     string    `json:"number"`
	Kind       string    `json:"kind"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
