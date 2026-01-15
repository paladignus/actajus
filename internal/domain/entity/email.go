// Package entity
package entity

import (
	"time"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type Email struct {
	IDEmails  uint
	Address   vo.Email
	DeletedAt *time.Time
}

func NewEmail(data dto.Email) Email {
	return Email{
		IDEmails: data.IDEmails,
		Address:  vo.Email(data.Address),
	}
}

func (e Email) Create() error {
	return e.validate()
}

func (e Email) Update() error {
	if e.IDEmails == 0 {
		return exception.ErrInvalidIDEmail
	}
	return e.validate()
}

func (e Email) IsDeleted() bool {
	return e.DeletedAt != nil
}

func (e Email) validate() error {
	if !e.Address.IsValid() {
		return exception.ErrInvalidEmail
	}
	return nil
}
