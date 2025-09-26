// Package dto
package dto

type AuthenticatedInput struct {
	CPF      string
	Password string
}

type AuthenticatedOutput struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	UserID       string   `json:"user_id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
}
