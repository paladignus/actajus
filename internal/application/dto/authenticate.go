// Package dto
package dto

type AuthenticateInput struct {
	CPF      string
	Password string
}

type Role struct {
	IDRole int    `json:"-"`
	Name   string `json:"name"`
}

type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type AuthenticateOutput struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	IDUser       string   `json:"id_user"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
}
