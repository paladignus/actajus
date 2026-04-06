// Package readmodel
package readmodel

import "time"

type CompanyReadModel struct {
	ID               int64
	Name             string
	TradeName        string
	CNPJ             string
	RegisteredByName string
	Address          *CompanyAddressReadModel
	Phones           []*CompanyPhoneReadModel
	Emails           []*CompanyEmailReadModel
	SocialMedia      []*CompanySocialMediaReadModel
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
