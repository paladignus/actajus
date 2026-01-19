// Package dto
package dto

type CreateCompanyRequest struct {
	Name         string `json:"name"`
	TradeName    string `json:"trade_name"`
	CNPJ         string `json:"cnpj"`
	RegisteredBy int    `json:"registered_by"`
}
