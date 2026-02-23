// Package dto
package dto

import "time"

type RequestPasswordResetReadModel struct {
	IDReset    int64     `json:"id_reset"`
	ResetToken string    `json:"reset_token"`
	ExpiresAt  time.Time `json:"expires_at"`
}
