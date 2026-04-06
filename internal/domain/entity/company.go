// Package entity
package entity

import (
	"time"

	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type Company struct {
	IDCompany    int
	RegisteredBy int
	Name         vo.Text
	TradeName    vo.Text
	CNPJ         vo.CNPJ
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func NewCompany(data readmodel.CompanyReadModel) Company {
	return Company{
		IDCompany:    data.IDCompany,
		RegisteredBy: data.RegisteredBy,
		Name:         vo.Text(data.Name),
		TradeName:    vo.Text(data.TradeName),
		CNPJ:         vo.CNPJ(data.CNPJ),
	}
}

func (c Company) Create() error {
	if c.RegisteredBy == 0 {
		return exception.ErrInvalidRegisteredBy
	}
	return c.validate()
}

func (c Company) Update() error {
	if c.IDCompany == 0 {
		return exception.ErrInvalidIDCompany
	}
	return c.validate()
}

func (c Company) IsDeleted() bool {
	return c.DeletedAt != nil
}

func (c Company) validate() error {
	if !c.Name.IsValid() {
		return exception.ErrInvalidName
	}
	if !c.TradeName.IsValid() {
		return exception.ErrInvalidTradeName
	}
	if !c.CNPJ.IsValid() {
		return exception.ErrInvalidCNPJ
	}
	return nil
}
