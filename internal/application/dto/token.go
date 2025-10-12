// Package dto
package dto

type TokenClaims struct {
	IDUser    string `json:"id_user"`
	TokenType string `json:"token_type"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
