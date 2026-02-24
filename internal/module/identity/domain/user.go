// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
)

type User struct {
	id           IDUser
	passwordHash PasswordHash
	isBlocked    bool
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    time.Time
}

type UserBuilder struct {
	u *User
}

func (b *UserBuilder) WithID(id int64) *UserBuilder {
	b.u.id = IDUser(id)
	return b
}

func (b *UserBuilder) WithPasswordHash(hash string) *UserBuilder {
	b.u.passwordHash = PasswordHash(hash)
	return b
}

func (b *UserBuilder) WithIsBlocked(isBlocked bool) *UserBuilder {
	b.u.isBlocked = isBlocked
	return b
}

func (b *UserBuilder) WithCreatedAt(t time.Time) *UserBuilder {
	b.u.createdAt = t
	return b
}

func (b *UserBuilder) WithUpdatedAt(t time.Time) *UserBuilder {
	b.u.updatedAt = t
	return b
}

func (b *UserBuilder) WithDeletedAt(t time.Time) *UserBuilder {
	b.u.deletedAt = t
	return b
}

func (b *UserBuilder) Build() (*User, error) {
	return b.u, nil
}

func (u *User) ID() IDUser                 { return u.id }
func (u *User) PasswordHash() PasswordHash { return u.passwordHash }
func (u *User) IsBlocked() bool            { return u.isBlocked }
func (u *User) CreatedAt() time.Time       { return u.createdAt }
func (u *User) UpdatedAt() time.Time       { return u.updatedAt }
func (u *User) DeletedAt() time.Time       { return u.deletedAt }

func (u *User) SetID(id int64) error {
	if u.id != 0 {
		return domain.NewFieldError("id", "user id is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "user id is invalid")
	}
	u.id = IDUser(id)
	return nil
}

func (u *User) SetPasswordHash(hash PasswordHash) {
	u.passwordHash = PasswordHash(hash.Value())
}
