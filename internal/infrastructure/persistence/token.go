// Package persistence
package persistence

import (
	"fmt"
	"time"

	"github.com/paladignus/actajus/internal/domain/repository"
)

type Token struct {
	repository repository.Cache
	// ShouldReturnError bool
	// ShouldReturnToken string
	// CallCount         int
}

func (t Token) SaveRevokedToken(token string, expiresAt time.Time) error {
	key := fmt.Sprintf("revoked_token:%s", token)
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}
	if err := t.repository.Set(key, "1", ttl); err != nil {
		return fmt.Errorf("failed to save revoked token: %w", err)
	}
	return nil
}

func (t Token) IsTokenRevoked(token string) (bool, error) {
	key := fmt.Sprintf("revoked_token:%s", token)
	exists, err := t.repository.Exists(key)
	if err != nil {
		return false, fmt.Errorf("failed to check if token is revoked: %w", err)
	}
	return exists > 0, nil
}

// func (t Token) ClearExpiredTokens() error {
// 	return nil
// }
