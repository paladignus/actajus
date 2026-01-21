// Package dto
package dto

import "github.com/paladignus/actajus/internal/module/address/application/dto"

type CompanyReadModel struct {
	ID           uint                  `json:"id"`
	Name         string                `json:"name"`
	TradeName    string                `json:"trade_name"`
	CNPJ         string                `json:"cnpj"`
	RegisteredBy uint                  `json:"registered_by"`
	Address      *dto.AddressReadModel `json:"address,omitzero"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
}
