// Package entity
package entity

import (
	"time"

	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type Enterprise struct {
	IDEnterprise int
	Name         vo.Text
	TradeName    vo.Text
	CNPJ         vo.CNPJ
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func NewEnterprise(name, tradeName, cnpj string) (Enterprise, error) {
	now := time.Now()
	e := Enterprise{
		Name:      vo.Text(name),
		TradeName: vo.Text(tradeName),
		CNPJ:      vo.CNPJ(cnpj),
		CreatedAt: now,
		UpdatedAt: now,
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
