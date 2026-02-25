// Package database
package database

import (
	"time"

	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	companyDTO "github.com/paladignus/actajus/internal/module/company/application/dto"
	emailDTO "github.com/paladignus/actajus/internal/module/email/application/dto"
	phoneDTO "github.com/paladignus/actajus/internal/module/phone/application/dto"
	socialMediaDTO "github.com/paladignus/actajus/internal/module/social_media/application/dto"
)

type companyScan struct {
	id               *uint
	registeredByName *string
	name             *string
	tradeName        *string
	cnpj             *string
	createdAt        *time.Time
	updatedAt        *time.Time
}

func (c companyScan) companyToDTO() *companyDTO.CompanyReadModel {
	if c.id == nil {
		return nil
	}
	registeredByName := ""
	if c.registeredByName != nil {
		registeredByName = *c.registeredByName
	}
	return &companyDTO.CompanyReadModel{
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
	id           *uint
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

func (a *addressScan) addressToDTO() *addrDTO.AddressReadModel {
	if a.id == nil {
		return nil
	}
	return &addrDTO.AddressReadModel{
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
	id        *uint
	address   *string
	createdAt *time.Time
	updatedAt *time.Time
}

func (e *emailScan) emailToDTO() *emailDTO.EmailReadModel {
	if e.id == nil {
		return nil
	}
	return &emailDTO.EmailReadModel{
		ID:        *e.id,
		Address:   *e.address,
		CreatedAt: *e.createdAt,
		UpdatedAt: *e.updatedAt,
	}
}

type phoneScan struct {
	id         *uint
	number     *string
	kind       *string
	department *string
	createdAt  *time.Time
	updatedAt  *time.Time
}

func (p *phoneScan) phoneToDTO() *phoneDTO.PhoneReadModel {
	if p.id == nil {
		return nil
	}
	return &phoneDTO.PhoneReadModel{
		ID:         *p.id,
		Number:     *p.number,
		Kind:       *p.kind,
		Department: *p.department,
		CreatedAt:  *p.createdAt,
		UpdatedAt:  *p.updatedAt,
	}
}

type socialMediaScan struct {
	id        *uint
	platform  *string
	url       *string
	createdAt *time.Time
	updatedAt *time.Time
}

func (sm socialMediaScan) socialMediaToDTO() *socialMediaDTO.SocialMediaReadModel {
	return &socialMediaDTO.SocialMediaReadModel{
		ID:        *sm.id,
		Platform:  *sm.platform,
		URL:       *sm.url,
		CreatedAt: *sm.createdAt,
		UpdatedAt: *sm.updatedAt,
	}
}
