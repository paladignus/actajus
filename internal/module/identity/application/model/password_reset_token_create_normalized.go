// Package model
package model

import (
	"time"
)

type PasswordResetTokenCreate struct {
	IDUser    int64
	Hash      [32]byte
	ExpiresAt time.Time
	CreatedAt time.Time
}
