package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokensut(t *testing.T) {
	sut := TokenClaims{
		IDUser: "user-id-123",
	}
	t.Run("should return the same ID user", func(t *testing.T) {
		assert.Equal(t, "user-id-123", sut.IDUser)
	})
	t.Run("should return an empty ID user", func(t *testing.T) {
		sut := TokenClaims{}
		assert.Equal(t, "", sut.IDUser)
	})
}

func TestTokenPair(t *testing.T) {
	sut := TokenPair{
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-456",
	}
	t.Run("should return the same access token and refresh token", func(t *testing.T) {
		assert.Equal(t, "access-token-123", sut.AccessToken)
		assert.Equal(t, "refresh-token-456", sut.RefreshToken)
	})
	t.Run("should return an empty access token and refresh token", func(t *testing.T) {
		sut := TokenPair{}
		assert.Equal(t, "", sut.AccessToken)
		assert.Equal(t, "", sut.RefreshToken)
	})
}
