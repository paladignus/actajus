// Package dto
package dto

import (
	addr "github.com/paladignus/actajus/internal/module/address/application/dto"
	email "github.com/paladignus/actajus/internal/module/email/application/dto"
	phone "github.com/paladignus/actajus/internal/module/phone/application/dto"
	socialMedia "github.com/paladignus/actajus/internal/module/social_media/application/dto"
)

type UpdateCompanyRequest struct {
	IDCompany   int64                                  `json:"id_company" validate:"required|numeric"`
	Name        string                                 `json:"name" validate:"required|min=3"`
	TradeName   string                                 `json:"trade_name" validate:"required|min=3"`
	CNPJ        string                                 `json:"cnpj" validate:"required|len=18|numeric"`
	Address     addr.UpdateAddressRequest              `json:"address,omitzero"`
	Phone       phone.UpdatePhoneRequest               `json:"phone,omitzero"`
	Email       email.UpdateEmailRequest               `json:"email,omitzero"`
	SocialMedia []socialMedia.UpdateSocialMediaRequest `json:"social_media,omitzero"`
}
