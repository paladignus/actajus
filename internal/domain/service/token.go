// Package service
package service

type AccessToken interface {
	GenereteAccessToken(string) (string, error)
	ValidateAccessToken(string) (string, error)
	RevokeAccessToken(string) error
}

type RefreshToken interface {
	GenereteRefreshToken(string) (string, error)
	ValidateRefreshToken(string) (string, error)
	RevokeRefreshToken(string) error
}

type TokenService interface {
	AccessToken
	RefreshToken
}
