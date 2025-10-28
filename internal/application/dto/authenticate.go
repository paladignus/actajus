// Package dto
package dto

type AuthenticateInput struct {
	CPF      string `json:"cpf"`
	Password string `json:"password"`
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

type Role struct {
	IDRole int    `json:"-"`
	Name   string `json:"name"`
}

type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type GetEmailByCPFInput struct {
	CPF string `json:"cpf"`
}

type GetEmailByCPFOutput struct {
	Email string `json:"email"`
}

type RecoverPasswordInput struct {
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RecoverPasswordOutput struct {
	RecoverToken string `json:"recover_token"`
}
