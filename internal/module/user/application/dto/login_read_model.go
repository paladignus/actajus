// Package dto
package dto

type LoginReadModel struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	IDSession        string `json:"id_session"`
	AccessExpiresAt  string `json:"access_expires_at"`
	RefreshExpiresAt string `json:"refresh_expires_at"`
}
