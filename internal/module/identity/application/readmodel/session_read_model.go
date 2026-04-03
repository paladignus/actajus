// Package readmodel
package readmodel

import "time"

type SessionReadModel struct {
	IDSession int64      `json:"id_session"`
	IDUser    int64      `json:"id_user"`
	UserEmail string     `json:"user_email"`
	IP        string     `json:"ip"`
	UserAgent string     `json:"user_agent"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	RotatedAt *time.Time `json:"rotated_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
