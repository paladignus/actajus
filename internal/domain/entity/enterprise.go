// Package entity
package entity

import (
	"time"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type Enterprise struct {
	IDEnterprise int
	RegisteredBy int
	Name         vo.Text
	TradeName    vo.Text
	CNPJ         vo.CNPJ
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func NewEnterprise(data dto.Enterprise) Enterprise {
	return Enterprise{
		IDEnterprise: data.IDEnterprise,
		RegisteredBy: data.RegisteredBy,
		Name:         vo.Text(data.Name),
		TradeName:    vo.Text(data.TradeName),
		CNPJ:         vo.CNPJ(data.CNPJ),
	}
}

func (e Enterprise) Create() error {
	if e.RegisteredBy == 0 {
		return exception.ErrInvalidRegisteredBy
	}
	return e.validate()
}

func (e Enterprise) Update() error {
	if e.IDEnterprise == 0 {
		return exception.ErrInvalidIDEnterprise
	}
	return e.validate()
}

func (e Enterprise) IsDeleted() bool {
	return e.DeletedAt != nil
}

func (e Enterprise) validate() error {
	if !e.Name.IsValid() {
		return exception.ErrInvalidName
	}
	if !e.TradeName.IsValid() {
		return exception.ErrInvalidTradeName
	}
	if !e.CNPJ.IsValid() {
		return exception.ErrInvalidCNPJ
	}
	return nil
}
