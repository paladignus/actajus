// Package domain
package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrUserBlocked        = errors.New("user_blocked")
	ErrSessionNotFound    = errors.New("session_not_found")
	ErrSessionRevoked     = errors.New("session_revoked")
	ErrSessionExpired     = errors.New("session_expired")
	ErrInvalidToken       = errors.New("invalid_token")
	ErrRefreshReuse       = errors.New("refresh_reuse_detected")
	ErrSessionLimit       = errors.New("session_limit_reached")
)
