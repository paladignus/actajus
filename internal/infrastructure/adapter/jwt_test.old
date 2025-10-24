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

func TestValidateRefreshToken(t *testing.T) {
	cfg := getTestConfig()
	adapter := NewJWTAdapter(cfg)
	validTokenPair, _ := adapter.GenerateTokenPair("user-456")
	expiredCfg := getTestConfig()
	expiredCfg.RefreshExpiry = -1 * time.Hour
	expiredAdapter := NewJWTAdapter(expiredCfg)
	expiredTokenPair, _ := expiredAdapter.GenerateTokenPair("user-456")
	wrongSecretCfg := getTestConfig()
	wrongSecretCfg.AccessSecret = "wrong-secret"
	wrongSecretToken, _ := NewJWTAdapter(wrongSecretCfg).GenerateTokenPair("user-123")
	invalidClaimsToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	invalidClaimsTokenString, _ := invalidClaimsToken.SignedString([]byte(cfg.RefreshSecret))
	tests := []struct {
		name     string
		token    string
		wantErr  error
		wantUser string
	}{
		{
			name:     "should return valid access token",
			token:    validTokenPair.RefreshToken,
			wantUser: "user-456",
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
			token:   expiredTokenPair.RefreshToken,
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
			claims, err := adapter.ValidateRefreshToken(tt.token)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantUser, claims.IDUser)
			assert.Equal(t, "refresh", claims.TokenType)
		})
	}
}

func TestRefreshAccessToken(t *testing.T) {
	cfg := getTestConfig()
	adapter := NewJWTAdapter(cfg)
	validTokenPair, _ := adapter.GenerateTokenPair("user-789")
	expiredCfg := getTestConfig()
	expiredCfg.RefreshExpiry = -1 * time.Hour
	expiredAdapter := NewJWTAdapter(expiredCfg)
	expiredTokenPair, _ := expiredAdapter.GenerateTokenPair("user-789")
	tests := []struct {
		name         string
		refreshToken string
		wantErr      bool
		validateUser string
	}{
		{
			name:         "should return an valid refresh token",
			refreshToken: validTokenPair.RefreshToken,
			validateUser: "user-789",
		},
		{
			name:         "should return an error to invalid token format",
			refreshToken: "invalid-token",
			wantErr:      true,
		},
		{
			name:         "should return an error to expired refresh token",
			refreshToken: expiredTokenPair.RefreshToken,
			wantErr:      true,
		},
		{
			name:         "should return an error to access token instead of refresh",
			refreshToken: validTokenPair.AccessToken,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTokenPair, err := adapter.RefreshAccessToken(tt.refreshToken)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, newTokenPair.AccessToken)
			assert.NotEmpty(t, newTokenPair.RefreshToken)
			accessClaims, err := adapter.ValidateAccessToken(newTokenPair.AccessToken)
			require.NoError(t, err)
			assert.Equal(t, tt.validateUser, accessClaims.IDUser)
			refreshClaims, err := adapter.ValidateRefreshToken(newTokenPair.RefreshToken)
			require.NoError(t, err)
			assert.Equal(t, tt.validateUser, refreshClaims.IDUser)
		})
	}
}
