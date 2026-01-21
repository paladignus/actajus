// Package dto
package dto

import "github.com/paladignus/actajus/internal/module/address/application/dto"

type CreateCompanyRequest struct {
	Name         string                   `json:"name"`
	TradeName    string                   `json:"trade_name"`
	CNPJ         string                   `json:"cnpj"`
	RegisteredBy uint                     `json:"registered_by"`
	Address      dto.CreateAddressRequest `json:"address,omitzero"`
}
