// Package dto
package dto

import "time"

type EmailReadModel struct {
	ID        int64     `json:"id"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
