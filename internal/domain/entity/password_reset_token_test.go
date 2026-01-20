// Package entity
package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasswordResetToken_IsValid(t *testing.T) {
	t.Run("should return an new PasswordResetToken", func(t *testing.T) {
		token := NewPasswordResetToken(1, "reset-token-789")
		assert.Equal(t, 1, token.IDUser)
		assert.Equal(t, "reset-token-789", token.Token)
		assert.NotEmpty(t, token.ExpiresAt)
		assert.NotEmpty(t, token.CreatedAt)
	})

	t.Run("should return false if the token is expired", func(t *testing.T) {
		token := PasswordResetToken{
			ExpiresAt: time.Now().Add(-time.Hour),
		}
		assert.False(t, token.IsValid())
	})

	t.Run("should return false if the token is used", func(t *testing.T) {
		now := time.Now()
		token := PasswordResetToken{
			UsedAt: &now,
		}
		assert.False(t, token.IsValid())
	})

	t.Run("should return true if the token is valid", func(t *testing.T) {
		token := PasswordResetToken{
			ExpiresAt: time.Now().Add(time.Hour),
		}
		assert.True(t, token.IsValid())
	})
}
