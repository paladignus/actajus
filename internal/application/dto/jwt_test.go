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

// func TestTokenRecover(t *testing.T) {
// 	sut := TokenRecover{
// 		ResetToken: "reset-token-789",
// 	}
// 	t.Run("should return the same reset token", func(t *testing.T) {
// 		assert.Equal(t, "reset-token-789", sut.ResetToken)
// 	})
// 	t.Run("should return an empty reset token", func(t *testing.T) {
// 		sut := TokenRecover{}
// 		assert.Equal(t, "", sut.ResetToken)
// 	})
// }

// func TestTokenStructsEmptyValues(t *testing.T) {
// 	emptysut := TokenClaims{}
// 	emptyPair := TokenPair{}
// 	emptyRecover := TokenRecover{}
// 	t.Run("should return empty values for sut", func(t *testing.T) {
// 		assert.Empty(t, emptysut.IDUser)
// 	})
// 	t.Run("should return empty values for TokenPair", func(t *testing.T) {
// 		assert.Empty(t, emptyPair.AccessToken)
// 		assert.Empty(t, emptyPair.RefreshToken)
// 	})
// 	t.Run("should return empty values for TokenRecover", func(t *testing.T) {
// 		assert.Empty(t, emptyRecover.ResetToken)
// 	})
// }
