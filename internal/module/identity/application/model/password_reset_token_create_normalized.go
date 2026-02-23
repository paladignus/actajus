// Package model
package model

import (
	"time"

	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type PasswordResetTokenCreate struct {
	IDUser    identity.IDUser
	Hash      [32]byte
	ExpiresAt time.Time
	CreatedAt time.Time
}
