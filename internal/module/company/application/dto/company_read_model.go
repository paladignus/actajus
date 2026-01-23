// Package dto
package dto

import (
	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	emailDTO "github.com/paladignus/actajus/internal/module/email/application/dto"
	phoneDTO "github.com/paladignus/actajus/internal/module/phone/application/dto"
	socialMediaDTO "github.com/paladignus/actajus/internal/module/social_media/application/dto"
)

type CompanyReadModel struct {
	ID           uint                                   `json:"id"`
	Name         string                                 `json:"name"`
	TradeName    string                                 `json:"trade_name"`
	CNPJ         string                                 `json:"cnpj"`
	RegisteredBy uint                                   `json:"registered_by"`
	Address      *addrDTO.AddressReadModel              `json:"address,omitzero"`
	Phone        *phoneDTO.PhoneReadModel               `json:"phone,omitzero"`
	Email        *emailDTO.EmailReadModel               `json:"email,omitzero"`
	SocialMedia  []*socialMediaDTO.SocialMediaReadModel `json:"social_media,omitzero"`
	CreatedAt    string                                 `json:"created_at"`
	UpdatedAt    string                                 `json:"updated_at"`
}
