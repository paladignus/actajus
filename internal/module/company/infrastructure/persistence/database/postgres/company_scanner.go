// Package postgres
package postgres

import (
	"time"

	companyReadModel "github.com/paladignus/actajus/internal/module/company/application/readmodel"
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

func (a *addressScan) addressToDTO() *companyReadModel.CompanyAddressReadModel {
	if a.id == nil {
		return nil
	}
	return &companyReadModel.CompanyAddressReadModel{
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

func (e *emailScan) emailToDTO() *companyReadModel.CompanyEmailReadModel {
	if e.id == nil {
		return nil
	}
	return &companyReadModel.CompanyEmailReadModel{
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

func (p *phoneScan) phoneToDTO() *companyReadModel.CompanyPhoneReadModel {
	if p.id == nil {
		return nil
	}
	return &companyReadModel.CompanyPhoneReadModel{
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

func (sm socialMediaScan) socialMediaToDTO() *companyReadModel.CompanySocialMediaReadModel {
	return &companyReadModel.CompanySocialMediaReadModel{
		ID:        *sm.id,
		Platform:  *sm.platform,
		URL:       *sm.url,
		CreatedAt: *sm.createdAt,
		UpdatedAt: *sm.updatedAt,
	}
}
