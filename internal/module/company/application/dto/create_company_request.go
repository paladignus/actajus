// Package dto
package dto

import (
	addr "github.com/paladignus/actajus/internal/module/address/application/dto"
	email "github.com/paladignus/actajus/internal/module/email/application/dto"
	phone "github.com/paladignus/actajus/internal/module/phone/application/dto"
)

type CreateCompanyRequest struct {
	Name         string                    `json:"name"`
	TradeName    string                    `json:"trade_name"`
	CNPJ         string                    `json:"cnpj"`
	RegisteredBy uint                      `json:"registered_by"`
	Address      addr.CreateAddressRequest `json:"address,omitzero"`
	Phone        phone.CreatePhoneRequest  `json:"phone,omitzero"`
	Email        email.CreateEmailRequest  `json:"email,omitzero"`
}
