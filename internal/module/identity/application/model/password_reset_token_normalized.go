// Package model
package model

import (
	"time"

	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type PasswordResetToken struct {
	ID        identity.IDPasswordReset
	IDUser    identity.IDUser
	Hash      [32]byte
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
