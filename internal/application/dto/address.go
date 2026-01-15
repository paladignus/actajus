// Package dto
package dto

type Address struct {
	IDAddress    uint   `json:"id_address"`
	Number       uint   `json:"number"`
	Zip          string `json:"zip"`
	Title        string `json:"title"`
	Street       string `json:"street"`
	Complement   string `json:"complement"`
	Reference    string `json:"reference"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country,omitzero"`
}
