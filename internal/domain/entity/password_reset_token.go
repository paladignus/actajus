// Package entity
package entity

import "time"

const TokenExpirationMinutes = 30

type PasswordResetToken struct {
	IDAccount int
	Token     string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func NewPasswordResetToken(IDAccount int, token string) PasswordResetToken {
	now := time.Now()
	return PasswordResetToken{
		IDAccount: IDAccount,
		Token:     token,
		ExpiresAt: now.Add(time.Minute * TokenExpirationMinutes),
		CreatedAt: now,
	}
}

func (p PasswordResetToken) IsExpired() bool {
	return p.ExpiresAt.Before(time.Now())
}

func (p PasswordResetToken) IsUsed() bool {
	return p.UsedAt != nil
}

func (p PasswordResetToken) IsValid() bool {
	return !p.IsExpired() && !p.IsUsed()
}

// func (p PasswordResetToken) MarkAsUsed() {
// 	now := time.Now()
// 	p.UsedAt = &now
// }
