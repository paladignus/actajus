// Package domain
package domain

import (
	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type User struct {
	id           vo.ID
	primaryEmail vo.Email
	passwordHash vo.PasswordHash
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

func (b *UserBuilder) WithID(id int64) *UserBuilder {
	b.u.id = vo.ID(id)
	return b
}

func (b *UserBuilder) WithPrimaryEmail(email string) *UserBuilder {
	b.u.primaryEmail = vo.Email(email)
	return b
}

func (b *UserBuilder) WithPasswordHash(hash string) *UserBuilder {
	b.u.passwordHash = vo.PasswordHash(hash)
	return b
}

func (b *UserBuilder) WithIsBlocked(blocked bool) *UserBuilder {
	b.u.isBlocked = blocked
	return b
}

func (b *UserBuilder) Build() (*User, error) {
	if b.u.primaryEmail.IsEmpty() && !b.u.primaryEmail.IsValid() {
		return nil, domain.NewFieldError("email", "email is invalid")
	}
	return b.u, nil
}

func (u *User) ID() vo.ID                     { return u.id }
func (u *User) PrimaryEmail() vo.Email        { return u.primaryEmail }
func (u *User) PasswordHash() vo.PasswordHash { return u.passwordHash }
func (u *User) IsBlocked() bool               { return u.isBlocked }

func (u *User) SetID(id vo.ID) error {
	if u.id != 0 {
		return domain.NewFieldError("id", "user id is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "user id is invalid")
	}
	u.id = id
	return nil
}

func (u *User) SetPasswordHash(hash string) {
	u.passwordHash = vo.PasswordHash(hash)
}
