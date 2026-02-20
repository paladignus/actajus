// Package dto
package dto

type RefreshTokenRequest struct {
	IDSession    string `json:"id_session"`
	RefreshToken string `json:"refresh_token"`
}
