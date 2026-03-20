// Package readmodel
package readmodel

import (
	"time"

	addrReadModel "github.com/paladignus/actajus/internal/module/address/application/readmodel"
	emailReadModel "github.com/paladignus/actajus/internal/module/email/application/readmodel"
	phoneReadModel "github.com/paladignus/actajus/internal/module/phone/application/readmodel"
	socialMediaReadModel "github.com/paladignus/actajus/internal/module/social_media/application/readmodel"
)

type CompanyReadModel struct {
	ID               int64
	Name             string
	TradeName        string
	CNPJ             string
	RegisteredByName string
	Addresses        *addrReadModel.AddressReadModel
	Phones           []*phoneReadModel.PhoneReadModel
	Emails           []*emailReadModel.EmailReadModel
	SocialMedia      []*socialMediaReadModel.SocialMediaReadModel
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
