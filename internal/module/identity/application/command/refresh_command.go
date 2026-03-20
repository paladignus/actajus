// Package dto
package dto

type RefreshCommand struct {
	IDSession    int64  `json:"id_session" validate:"required|min=1"`
	RefreshToken string `json:"refresh_token" validate:"required|min=10"`
	IP           string `json:"ip"`
	UserAgent    string `json:"user_agent"`
}
