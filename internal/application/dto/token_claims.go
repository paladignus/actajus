// Package dto
package dto

type TokenClaims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
}
