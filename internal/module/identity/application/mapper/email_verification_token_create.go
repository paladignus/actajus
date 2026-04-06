// Package mapper
package mapper

import "time"

type EmailVerificationTokenCreate struct {
	IDUser    int64
	IDEmail   int64
	Hash      [32]byte
	ExpiresAt time.Time
	CreatedAt time.Time
}
