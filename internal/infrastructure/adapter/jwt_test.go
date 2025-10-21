// Package adapter
package adapter

import (
	"errors"
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

type faultyReader struct{}

func (f faultyReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("random read error")
}

func TestJWTAdapter_GenerateTokenPair(t *testing.T) {
	cfg := config.JWTConfig{
		AccessSecret:  "access-secret-key",
		RefreshSecret: "refresh-secret-key",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}
	adapter := NewJWTAdapter(cfg)
	tests := []struct {
		name      string
		userID    string
		wantError bool
	}{
		{
			name:      "should generate token pair",
			userID:    "user-123",
			wantError: false,
		},
		{
			name:      "should generate token pair with different user ID",
			userID:    "user-456",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenPair, err := adapter.GenerateTokenPair(tt.userID)
			if tt.wantError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, tokenPair.AccessToken)
			assert.NotEmpty(t, tokenPair.RefreshToken)
			assert.True(t, len(tokenPair.AccessToken) > 0)
			assert.True(t, len(tokenPair.RefreshToken) > 0)
			// Validate that tokens are different
			assert.NotEqual(t, tokenPair.AccessToken, tokenPair.RefreshToken)
			// Validate access token structure
			accessClaims, err := adapter.ValidateAccessToken(tokenPair.AccessToken)
			assert.NoError(t, err)
			assert.Equal(t, tt.userID, accessClaims.IDUser)
			assert.Equal(t, "access", accessClaims.TokenType)
			// Validate refresh token structure
			refreshClaims, err := adapter.ValidateRefreshToken(tokenPair.RefreshToken)
			assert.NoError(t, err)
			assert.Equal(t, tt.userID, refreshClaims.IDUser)
			assert.Equal(t, "refresh", refreshClaims.TokenType)
		})
	}
}

// func TestJWT(t *testing.T) {
// 	config := config.Load()
// 	sut := spy.NewJWTAdapter(config.JWT.AccessSecret, config.JWT.AccessExpiry)
// 	t.Run("should create a new JWT adapter instance", func(t *testing.T) {
// 		if sut.SecretKey != config.JWT.AccessSecret || sut.ExpiresIn != config.JWT.AccessExpiry {
// 			t.Errorf("expected SecretKey %s and ExpiresIn %v, got SecretKey %s and ExpiresIn %v", sut.SecretKey, sut.ExpiresIn, config.JWT.AccessSecret, config.JWT.AccessExpiry)
// 		}
// 	})
// 	t.Run("should have correct expiration time", func(t *testing.T) {
// 		expectedExpiry := 15 * time.Minute
// 		if sut.ExpiresIn != expectedExpiry {
// 			t.Errorf("expected ExpiresIn %v, got %v", expectedExpiry, sut.ExpiresIn)
// 		}
// 	})
// 	t.Run("should call NewToken once", func(t *testing.T) {
// 		sut.NewToken()
// 		if sut.CallsCount != 1 {
// 			t.Errorf("expected calCount to be 1, got %d", sut.CallsCount)
// 		}
// 	})
// 	t.Run("should return error when ErrGenerateToken is set", func(t *testing.T) {
// 		sut.ErrGenerateToken = errors.New("error generating token")
// 		_, err := sut.NewToken()
// 		if err == nil {
// 			t.Error("expected error to be not nil")
// 		}
// 	})
// 	t.Run("should return valid token when no error", func(t *testing.T) {
// 		sut.ErrGenerateToken = nil
// 		token, err := sut.NewToken()
// 		if err != nil {
// 			t.Error("expected error to be nil, got", err)
// 		}
// 		if token != "valid_token" {
// 			t.Errorf("expected token to be 'valid_token', got %s", token)
// 		}
// 	})
// }
