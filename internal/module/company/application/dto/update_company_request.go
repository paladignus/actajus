// Package dto
package dto

type UpdateCompanyRequest struct {
	Name      string `json:"name"`
	TradeName string `json:"trade_name"`
	CNPJ      string `json:"cnpj"`
}
