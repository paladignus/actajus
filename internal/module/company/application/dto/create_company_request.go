// Package dto
package dto

import (
	addr "github.com/paladignus/actajus/internal/module/address/application/dto"
	email "github.com/paladignus/actajus/internal/module/email/application/dto"
	phone "github.com/paladignus/actajus/internal/module/phone/application/dto"
	socialMedia "github.com/paladignus/actajus/internal/module/social_media/application/dto"
)

type CreateCompanyRequest struct {
	Name         string                                 `json:"name" validate:"required|min=3"`
	TradeName    string                                 `json:"trade_name" validate:"required|min=3"`
	CNPJ         string                                 `json:"cnpj" validate:"required|len=18|numeric"`
	RegisteredBy uint                                   `json:"registered_by" validate:"required|min=1"`
	Address      addr.CreateAddressRequest              `json:"address,omitzero"`
	Phone        phone.CreatePhoneRequest               `json:"phone,omitzero"`
	Email        email.CreateEmailRequest               `json:"email,omitzero"`
	SocialMedia  []socialMedia.CreateSocialMediaRequest `json:"social_media,omitzero"`
}
