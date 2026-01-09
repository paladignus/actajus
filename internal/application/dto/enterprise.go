// Package dto
package dto

type EnterpriseInput struct {
	RegisteredBy int    `json:"registered_by"`
	Name         string `json:"name"`
	TradeName    string `json:"trade_name"`
	CNPJ         string `json:"cnpj"`
}
