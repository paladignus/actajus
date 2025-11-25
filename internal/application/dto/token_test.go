package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenClaims(t *testing.T) {
	claims := TokenClaims{
		IDUser: "user-id-123",
	}
	t.Run("should return the same ID user", func(t *testing.T) {
		assert.Equal(t, "user-id-123", claims.IDUser)
	})
	t.Run("should return an empty ID user", func(t *testing.T) {
		claims := TokenClaims{}
		assert.Equal(t, "", claims.IDUser)
	})
}

func TestTokenPair(t *testing.T) {
	pair := TokenPair{
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-456",
	}
	t.Run("should return the same access token and refresh token", func(t *testing.T) {
		assert.Equal(t, "access-token-123", pair.AccessToken)
		assert.Equal(t, "refresh-token-456", pair.RefreshToken)
	})
	t.Run("should return an empty access token and refresh token", func(t *testing.T) {
		pair := TokenPair{}
		assert.Equal(t, "", pair.AccessToken)
		assert.Equal(t, "", pair.RefreshToken)
	})
}

func TestTokenRecover(t *testing.T) {
	recoverToken := TokenRecover{
		ResetToken: "reset-token-789",
	}
	t.Run("should return the same reset token", func(t *testing.T) {
		assert.Equal(t, "reset-token-789", recoverToken.ResetToken)
	})
	t.Run("should return an empty reset token", func(t *testing.T) {
		recoverToken := TokenRecover{}
		assert.Equal(t, "", recoverToken.ResetToken)
	})
}

func TestTokenStructsEmptyValues(t *testing.T) {
	emptyClaims := TokenClaims{}
	emptyPair := TokenPair{}
	emptyRecover := TokenRecover{}
	t.Run("should return empty values for claims", func(t *testing.T) {
		assert.Empty(t, emptyClaims.IDUser)
	})
	t.Run("should return empty values for TokenPair", func(t *testing.T) {
		assert.Empty(t, emptyPair.AccessToken)
		assert.Empty(t, emptyPair.RefreshToken)
	})
	t.Run("should return empty values for TokenRecover", func(t *testing.T) {
		assert.Empty(t, emptyRecover.ResetToken)
	})
}

