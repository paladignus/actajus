// Package domain
package domain

import (
	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type User struct {
	id           IDUser
	primaryEmail vo.Email
	passwordHash PasswordHash
	isBlocked    bool
}

type UserBuilder struct {
	u *User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		&User{},
	}
}

func (b *UserBuilder) WithID(id IDUser) *UserBuilder {
	b.u.id = id
	return b
}

func (b *UserBuilder) WithPrimaryEmail(email vo.Email) *UserBuilder {
	b.u.primaryEmail = email
	return b
}

func (b *UserBuilder) WithPasswordHash(hash PasswordHash) *UserBuilder {
	b.u.passwordHash = hash
	return b
}

func (b *UserBuilder) WithIsBlocked(blocked bool) *UserBuilder {
	b.u.isBlocked = blocked
	return b
}

func (b *UserBuilder) Build() (*User, error) {
	return b.u, nil
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

func (u *User) SetID(id IDUser) error {
	if u.id != 0 {
		return domain.NewFieldError("id", "user id is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "user id is invalid")
	}
	u.id = id
	return nil
}

func (u *User) SetPasswordHash(hash PasswordHash) {
	u.passwordHash = PasswordHash(hash.Value())
}
