// Package adapter
package adapter

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/paladignus/actajus/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestConfig() config.JWTConfig {
	return config.JWTConfig{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}
}

func TestGenerateTokenPair(t *testing.T) {
	adapter := NewJWTAdapter(getTestConfig())
	tokenPair, err := adapter.GenerateTokenPair("user-123")
	require.NoError(t, err)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.NotEqual(t, tokenPair.AccessToken, tokenPair.RefreshToken)
}

func TestValidateAccessToken(t *testing.T) {
	cfg := getTestConfig()
	adapter := NewJWTAdapter(cfg)
	validTokenPair, _ := adapter.GenerateTokenPair("user-123")
	expiredCfg := getTestConfig()
	expiredCfg.AccessExpiry = -1 * time.Hour
	expiredAdapter := NewJWTAdapter(expiredCfg)
	expiredTokenPair, _ := expiredAdapter.GenerateTokenPair("user-123")
	wrongSecretCfg := getTestConfig()
	wrongSecretCfg.AccessSecret = "wrong-secret"
	wrongSecretToken, _ := NewJWTAdapter(wrongSecretCfg).GenerateTokenPair("user-123")
	invalidClaimsToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	invalidClaimsTokenString, _ := invalidClaimsToken.SignedString([]byte(cfg.AccessSecret))
	tests := []struct {
		name     string
		token    string
		wantErr  error
		wantUser string
	}{
		{
			name:     "should return valid access token",
			token:    validTokenPair.AccessToken,
			wantUser: "user-123",
		},
		{
			name:    "should return invalid token error to token invalid",
			token:   "invalid-token",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "should return invalid token error to empty token",
			token:   "",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "should return expired token error",
			token:   expiredTokenPair.AccessToken,
			wantErr: ErrExpiredToken,
		},
		{
			name:    "should return invalid token error to wrong secret token",
			token:   wrongSecretToken.AccessToken,
			wantErr: ErrInvalidToken,
		},
		{
			name:    "should return invalid token type error",
			token:   invalidClaimsTokenString,
			wantErr: ErrInvalidTokenType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := adapter.ValidateAccessToken(tt.token)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantUser, claims.IDUser)
			assert.Equal(t, "access", claims.TokenType)
		})
	}
}
