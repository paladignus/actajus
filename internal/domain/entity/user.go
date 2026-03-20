// Package entity
package entity

import (
	"time"

	"github.com/paladignus/actajus/internal/application/command"
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

func NewUser(input command.CreateUserCommand) (User, error) {
	a := User{
		Password:    vo.Password(input.Password),
		Avatar:      vo.File(input.Avatar),
		LastLoginAt: input.LastLoginAt,
	}
	if !a.Password.IsValid() {
		return a, exception.ErrInvalidPassword
	}
	if !a.Avatar.IsValid() {
		return a, exception.ErrInvalidAvatar
	}
	now := time.Now()
	if a.LastLoginAt.IsZero() || a.LastLoginAt.After(now) {
		return a, exception.ErrInvalidLastLoginAt
	}
	return a, nil
}

func (u User) IsDeleted() bool {
	return u.DeletedAt != nil
}
