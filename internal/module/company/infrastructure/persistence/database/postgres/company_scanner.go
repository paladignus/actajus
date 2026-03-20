// Package postgres
package postgres

import (
	"time"

	addrReadModel "github.com/paladignus/actajus/internal/module/address/application/readmodel"
	companyReadModel "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	emailReadModel "github.com/paladignus/actajus/internal/module/email/application/readmodel"
	phoneReadModel "github.com/paladignus/actajus/internal/module/phone/application/readmodel"
	socialMediaReadModel "github.com/paladignus/actajus/internal/module/social_media/application/readmodel"
)

type companyScan struct {
	id               *int64
	registeredByName *string
	name             *string
	tradeName        *string
	cnpj             *string
	createdAt        *time.Time
	updatedAt        *time.Time
}

func (c companyScan) companyToDTO() *companyReadModel.CompanyReadModel {
	if c.id == nil {
		return nil
	}
	registeredByName := ""
	if c.registeredByName != nil {
		registeredByName = *c.registeredByName
	}
	return &companyReadModel.CompanyReadModel{
		ID:               *c.id,
		Name:             *c.name,
		RegisteredByName: registeredByName,
		TradeName:        *c.tradeName,
		CNPJ:             *c.cnpj,
		CreatedAt:        *c.createdAt,
		UpdatedAt:        *c.updatedAt,
	}
}

// interno ao repositório, não exportado
type addressScan struct {
	id           *int64
	zip          *string
	title        *string
	street       *string
	complement   *string
	reference    *string
	number       *uint
	neighborhood *string
	city         *string
	state        *string
	country      *string
	createdAt    *time.Time
	updatedAt    *time.Time
}

func (a *addressScan) addressToDTO() *addrReadModel.AddressReadModel {
	if a.id == nil {
		return nil
	}
	return &addrReadModel.AddressReadModel{
		ID:           *a.id,
		ZIP:          *a.zip,
		Title:        *a.title,
		Street:       *a.street,
		Number:       *a.number,
		Complement:   a.complement, // já é *string, passa direto
		Reference:    a.reference,
		Neighborhood: *a.neighborhood,
		City:         *a.city,
		State:        *a.state,
		Country:      *a.country,
		CreatedAt:    *a.createdAt,
		UpdatedAt:    *a.updatedAt,
	}
}

type emailScan struct {
	id        *int64
	address   *string
	createdAt *time.Time
	updatedAt *time.Time
}

func (e *emailScan) emailToDTO() *emailReadModel.EmailReadModel {
	if e.id == nil {
		return nil
	}
	return &emailReadModel.EmailReadModel{
		ID:        *e.id,
		Address:   *e.address,
		CreatedAt: *e.createdAt,
		UpdatedAt: *e.updatedAt,
	}
}

type phoneScan struct {
	id         *int64
	number     *string
	kind       *string
	department *string
	createdAt  *time.Time
	updatedAt  *time.Time
}

func (p *phoneScan) phoneToDTO() *phoneReadModel.PhoneReadModel {
	if p.id == nil {
		return nil
	}
	return &phoneReadModel.PhoneReadModel{
		ID:         *p.id,
		Number:     *p.number,
		Kind:       *p.kind,
		Department: *p.department,
		CreatedAt:  *p.createdAt,
		UpdatedAt:  *p.updatedAt,
	}
}

type socialMediaScan struct {
	id        *int64
	platform  *string
	url       *string
	createdAt *time.Time
	updatedAt *time.Time
}

func (sm socialMediaScan) socialMediaToDTO() *socialMediaReadModel.SocialMediaReadModel {
	return &socialMediaReadModel.SocialMediaReadModel{
		ID:        *sm.id,
		Platform:  *sm.platform,
		URL:       *sm.url,
		CreatedAt: *sm.createdAt,
		UpdatedAt: *sm.updatedAt,
	}
}
