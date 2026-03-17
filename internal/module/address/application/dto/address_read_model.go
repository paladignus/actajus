// Package dto
package dto

import "time"

type AddressReadModel struct {
	ID           int64     `json:"id"`
	ZIP          string    `json:"zip"`
	Title        string    `json:"title"`
	Street       string    `json:"street"`
	Number       uint      `json:"number"`
	Complement   *string   `json:"complement,omitzero"`
	Reference    *string   `json:"reference,omitzero"`
	Neighborhood string    `json:"neighborhood"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Country      string    `json:"country"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
