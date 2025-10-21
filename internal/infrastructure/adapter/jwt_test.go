// Package adapter
package adapter

import (
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/application/dto"
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

func TestJWTAdapter_GenerateTokenPair_Error(t *testing.T) {
	t.Run("error - random generator failure", func(t *testing.T) {
		originalReader := rand.Reader
		rand.Reader = faultyReader{}
		defer func() { rand.Reader = originalReader }()
		cfg := config.JWTConfig{
			AccessSecret:  "access-secret-key",
			RefreshSecret: "refresh-secret-key",
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
		}
		adapter := NewJWTAdapter(cfg).(jwtAdapter)
		tokenPair, err := adapter.GenerateTokenPair("user-123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to generate access token")
		assert.Equal(t, dto.TokenPair{}, tokenPair)
	})
}
