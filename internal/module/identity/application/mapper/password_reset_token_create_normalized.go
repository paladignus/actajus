// Package mapper
package mapper

import (
	"time"
)

type PasswordResetTokenCreate struct {
	IDUser    int64
	Hash      [32]byte
	ExpiresAt time.Time
	CreatedAt time.Time
}
