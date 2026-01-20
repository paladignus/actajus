// Package dto
package dto

import addressDTO "github.com/paladignus/actajus/internal/module/address/application/dto"

type CompanyResponse struct {
	ID           uint                        `json:"id"`
	Name         string                      `json:"name"`
	TradeName    string                      `json:"trade_name"`
	CNPJ         string                      `json:"cnpj"`
	RegisteredBy uint                        `json:"registered_by"`
	Address      *addressDTO.AddressResponse `json:"address,omitempty"`
	CreatedAt    string                      `json:"created_at"`
	UpdatedAt    string                      `json:"updated_at"`
}

type CompanyListResponse struct {
	Companies []CompanyResponse `json:"companies"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	Total     int64             `json:"total"`
}
