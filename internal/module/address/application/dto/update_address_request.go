// Package dto
package dto

type UpdateAddressRequest struct {
	IDAddress    int64  `json:"id_address" validate:"required|numeric"`
	ZIP          string `json:"zip" validate:"required|len=9|numeric"`
	Title        string `json:"title" validate:"required|min=3|max=100"`
	Street       string `json:"street" validate:"required|min=3|max=100"`
	Number       uint   `json:"number" validate:"required|min=1"`
	Complement   string `json:"complement,omitzero"`
	Reference    string `json:"reference,omitzero"`
	Neighborhood string `json:"neighborhood" validate:"required|min=3|max=100"`
	City         string `json:"city" validate:"required|min=2|max=100"`
	State        string `json:"state" validate:"required|max=100"`
	Country      string `json:"country" validate:"required|min=2|max=100"`
}
