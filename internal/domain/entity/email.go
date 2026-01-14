// Package entity
package entity

import (
	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type Email struct {
	IDEmails uint
	Address  vo.Email
}

func NewEmail(input dto.Email) (Email, error) {
	e := Email{
		Address: vo.Email(input.Address),
	}
	if !e.Address.IsValid() {
		return e, exception.ErrInvalidEmail
	}
	return e, nil
}
