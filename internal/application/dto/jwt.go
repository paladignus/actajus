// Package dto
package dto

type TokenClaims struct {
	IDUser string `json:"id_user"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// type TokenRecover struct {
// 	ResetToken string `json:"recover_token"`
// }
