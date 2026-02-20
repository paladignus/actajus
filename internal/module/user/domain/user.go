// Package domain
package domain

import (
	"time"

	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type User struct {
	id        vo.ID
	password  vo.Password
	avatar    vo.File
	isBlocked bool
	roles     []string
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

type UserBuilder struct {
	user *User
}

func NewUserBuilder() *UserBuilder {
	now := time.Now()
	return &UserBuilder{
		&User{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (u *UserBuilder) WithID(id int64) *User {
	u.user.id = vo.ID(id)
	return u.user
}

func (u *UserBuilder) WithPassword(password string) *User {
	u.user.password = vo.Password(password)
	return u.user
}

func (u *UserBuilder) WithAvatar(avatar string) *User {
	u.user.avatar = vo.File(avatar)
	return u.user
}

func (u *UserBuilder) WithIsBlocked(status bool) *User {
	u.user.isBlocked = status
	return u.user
}

func (u *UserBuilder) WithRoles(roles []string) *User {
	u.user.roles = roles
	return u.user
}

func (u *UserBuilder) WithCreatedAt(createdAt time.Time) *User {
	u.user.createdAt = createdAt
	return u.user
}

func (u *UserBuilder) WithUpdatedAt(updatedAt time.Time) *User {
	u.user.updatedAt = updatedAt
	return u.user
}

func (u *UserBuilder) WithDeletedAt(deletedAt *time.Time) *User {
	u.user.deletedAt = deletedAt
	return u.user
}

func (u *UserBuilder) Build() (*User, error) {
	return u.user, nil
}

func (u *User) ID() vo.ID             { return u.id }
func (u *User) Password() vo.Password { return u.password }
func (u *User) Avatar() vo.File       { return u.avatar }
func (u *User) IsBlocked() bool       { return u.isBlocked }
func (u *User) Roles() []string       { return u.roles }
func (u *User) CreatedAt() time.Time  { return u.createdAt }
func (u *User) UpdatedAt() time.Time  { return u.updatedAt }
func (u *User) DeletedAt() *time.Time { return u.deletedAt }
