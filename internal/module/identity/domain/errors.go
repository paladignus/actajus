// Package domain
package domain

import "errors"

var (
	ErrInvalidCredentials        = errors.New("invalid_credentials")
	ErrUserBlocked               = errors.New("user_blocked")
	ErrEmailNotVerified          = errors.New("email_not_verified")
	ErrSessionNotFound           = errors.New("session_not_found")
	ErrSessionRevoked            = errors.New("session_revoked")
	ErrSessionExpired            = errors.New("session_expired")
	ErrInvalidToken              = errors.New("invalid_token")
	ErrRefreshReuse              = errors.New("refresh_reuse_detected")
	ErrSessionLimit              = errors.New("session_limit_reached")
	ErrResetTokenNotFound        = errors.New("reset_token_not_found")
	ErrResetTokenExpired         = errors.New("reset_token_expired")
	ErrResetTokenUsed            = errors.New("reset_token_used")
	ErrResetTokenInvalid         = errors.New("reset_token_invalid")
	ErrEmailAlreadyExists        = errors.New("email_already_exists")
	ErrEmailAlreadyVerified      = errors.New("email_already_verified")
	ErrEmailVerificationNotFound = errors.New("email_verification_not_found")
	ErrEmailVerificationExpired  = errors.New("email_verification_expired")
	ErrEmailVerificationUsed     = errors.New("email_verification_used")
	ErrEmailVerificationInvalid  = errors.New("email_verification_invalid")
	ErrMissingAccessToken        = errors.New("missing_access_token")
	ErrSessionNotActive          = errors.New("session_not_active")
)
