// Package mapper
package mapper

import "time"

type EmailVerificationToken struct {
	ID        int64
	IDUser    int64
	IDEmail   int64
	Hash      [32]byte
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
