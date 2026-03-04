// Package model
package model

import (
	"time"
)

type PasswordResetToken struct {
	ID        int64
	IDUser    int64
	Hash      [32]byte
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
