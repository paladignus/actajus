// Package entity
package entity

import (
	"time"

	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/paladignus/actajus/internal/domain/exception"
)

type Phone struct {
	IDPhone    int
	Number     string
	Kind       string
	Department string
	DeletedAt  *time.Time
}

func NewPhone(phone readmodel.PhoneReadModel) Phone {
	return Phone{
		IDPhone:    phone.IDPhone,
		Number:     phone.Number,
		Kind:       phone.Kind,
		Department: phone.Department,
	}
}

func (e Phone) Update() error {
	if e.IDPhone == 0 {
		return exception.ErrInvalidIDPhone
	}
	return nil
}
