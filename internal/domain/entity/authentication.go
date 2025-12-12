// Package entity
package entity

import (
	"time"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type Authentication struct {
	IDPerson    int
	Password    vo.Password
	Avatar      vo.File
	LastLoginAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func NewAuthentication(authentication dto.AuthenticationInput) (Authentication, error) {
	a := Authentication{
		IDPerson:    authentication.IDPeople,
		Password:    vo.Password(authentication.Password),
		Avatar:      vo.File(authentication.Avatar),
		LastLoginAt: authentication.LastLoginAt,
		CreatedAt:   authentication.CreatedAt,
		UpdatedAt:   authentication.UpdatedAt,
		DeletedAt:   authentication.DeletedAt,
	}
	if !a.Password.IsValid() {
		return a, exception.ErrInvalidPassword
	}
	if !a.Avatar.IsValid() {
		return a, exception.ErrInvalidAvatar
	}
	now := time.Now()
	if !a.LastLoginAt.IsZero() || a.LastLoginAt.Before(now) {
		return a, exception.ErrInvalidLastLoginAt
	}
	return a, nil
}

func (a Authentication) IsDeleted() bool {
	return a.DeletedAt != nil
}
