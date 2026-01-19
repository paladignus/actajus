// Package dto
package dto

type CompanyResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	TradeName    string `json:"trade_name"`
	CNPJ         string `json:"cnpj"`
	RegisteredBy int    `json:"registered_by"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CompanyListResponse struct {
	Companies []CompanyResponse `json:"companies"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	Total     int64             `json:"total"`
}
