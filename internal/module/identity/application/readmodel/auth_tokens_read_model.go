// Package dto
package readmodel

import "time"

type AuthTokensReadModel struct {
	IDSession        int64     `json:"id_session"`
	IDUser           int64     `json:"id_user"`
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}
