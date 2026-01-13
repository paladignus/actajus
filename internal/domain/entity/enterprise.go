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

func NewEnterprise(enterprise dto.Enterprise) (Enterprise, error) {
	now := time.Now()
	e := Enterprise{
		RegisteredBy: enterprise.RegisteredBy,
		Name:         vo.Text(enterprise.Name),
		TradeName:    vo.Text(enterprise.TradeName),
		CNPJ:         vo.CNPJ(enterprise.CNPJ),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if e.RegisteredBy == 0 {
		return e, exception.ErrInvalidRegisteredBy
	}
	if !e.Name.IsValid() {
		return e, exception.ErrInvalidName
	}
	if !e.TradeName.IsValid() {
		return e, exception.ErrInvalidTradeName
	}
	if !e.CNPJ.IsValid() {
		return e, exception.ErrInvalidCNPJ
	}
	return e, nil
}

func (e Enterprise) IsDeleted() bool {
	return e.DeletedAt != nil
}
