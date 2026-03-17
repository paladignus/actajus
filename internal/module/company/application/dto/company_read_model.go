// Package dto
package dto

import (
	"time"

	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	emailDTO "github.com/paladignus/actajus/internal/module/email/application/dto"
	phoneDTO "github.com/paladignus/actajus/internal/module/phone/application/dto"
	socialMediaDTO "github.com/paladignus/actajus/internal/module/social_media/application/dto"
)

type CompanyReadModel struct {
	ID               int64                                  `json:"id"`
	Name             string                                 `json:"name"`
	TradeName        string                                 `json:"trade_name"`
	CNPJ             string                                 `json:"cnpj"`
	RegisteredByName string                                 `json:"registered_by,omitzero"`
	Addresses        *addrDTO.AddressReadModel              `json:"address,omitzero"`
	Phones           []*phoneDTO.PhoneReadModel             `json:"phone,omitzero"`
	Emails           []*emailDTO.EmailReadModel             `json:"email,omitzero"`
	SocialMedia      []*socialMediaDTO.SocialMediaReadModel `json:"social_media,omitzero"`
	CreatedAt        time.Time                              `json:"created_at"`
	UpdatedAt        time.Time                              `json:"updated_at"`
}
