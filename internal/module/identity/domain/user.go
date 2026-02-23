// Package domain
package domain

import (
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type User struct {
	id           IDUser
	primaryEmail vo.Email
	passwordHash PasswordHash
	isBlocked    bool
}

func NewUser(id IDUser, primaryEmail vo.Email, hash PasswordHash, blocked bool) (*User, error) {
	return &User{
		id:           id,
		primaryEmail: primaryEmail,
		passwordHash: hash,
		isBlocked:    blocked,
	}, nil
}

func (u *User) ID() IDUser                 { return u.id }
func (u *User) PrimaryEmail() vo.Email     { return u.primaryEmail }
func (u *User) PasswordHash() PasswordHash { return u.passwordHash }
func (u *User) IsBlocked() bool            { return u.isBlocked }
