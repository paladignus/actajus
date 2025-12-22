// Package entity
package entity

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Run("should return an user", func(t *testing.T) {
		now := time.Now()
		user, err := NewUser(dto.UserInput{
			Password:    "P@ssword123",
			Avatar:      "avatar.png",
			LastLoginAt: now,
		})
		assert.NoError(t, err)
		assert.Equal(t, "P@ssword123", user.Password.Value())
		assert.Equal(t, "avatar.png", user.Avatar.Value())
		assert.Equal(t, now, user.LastLoginAt)
	})

	t.Run("should return false if the password is invalid", func(t *testing.T) {
		_, err := NewUser(dto.UserInput{})
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidPassword, err)
	})

	t.Run("should return false if the avatar is invalid", func(t *testing.T) {
		_, err := NewUser(dto.UserInput{
			Password: "P@ssword123",
		})
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidAvatar, err)
	})

	t.Run("should return false if the last login at is invalid", func(t *testing.T) {
		_, err := NewUser(dto.UserInput{
			Password:    "P@ssword123",
			Avatar:      "avatar.png",
			LastLoginAt: time.Now().Add(time.Hour * 24),
		})
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidLastLoginAt, err)
	})

	t.Run("should return false if the user is not deleted", func(t *testing.T) {
		user := User{}
		assert.False(t, user.IsDeleted())
	})
}
