// Package entity
package entity

import (
	"time"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type User struct {
	IDUser      int
	Password    vo.Password
	Avatar      vo.File
	LastLoginAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func NewUser(user dto.UserInput) (User, error) {
	a := User{
		// IDPerson:    user.IDPeople,
		Password:    vo.Password(user.Password),
		Avatar:      vo.File(user.Avatar),
		LastLoginAt: user.LastLoginAt,
		// CreatedAt:   user.CreatedAt,
		// UpdatedAt:   user.UpdatedAt,
		// DeletedAt:   user.DeletedAt,
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

func (u User) IsDeleted() bool {
	return u.DeletedAt != nil
}
