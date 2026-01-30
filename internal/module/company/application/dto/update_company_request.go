// Package dto
package dto

import (
	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
)

type UpdateCompanyRequest struct {
	IDCompany uint                         `json:"id_company" validate:"required|numeric"`
	Name      string                       `json:"name" validate:"required|min=3"`
	TradeName string                       `json:"trade_name" validate:"required|min=3"`
	CNPJ      string                       `json:"cnpj" validate:"required|len=18|numeric"`
	Address   addrDTO.UpdateAddressRequest `json:"address,omitzero"`
}
