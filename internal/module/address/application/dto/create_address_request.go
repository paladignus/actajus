// Package dto
package dto

type CreateAddressRequest struct {
	ZIP          string `json:"zip" validate:"required|len=8|numeric"`
	Title        string `json:"title"`
	Street       string `json:"street"`
	Number       uint   `json:"number"`
	Complement   string `json:"complement,omitzero"`
	Reference    string `json:"reference,omitzero"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
}
