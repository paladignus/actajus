// Package dto
package dto

type Enterprise struct {
	RegisteredBy int    `json:"registered_by"`
	Name         string `json:"name"`
	TradeName    string `json:"trade_name"`
	CNPJ         string `json:"cnpj"`
}

type EnterpriseInput struct {
	Enterprise `json:"enterprise"`
	Address    `json:"address"`
}
