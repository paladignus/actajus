// Package command
package command

type LoginCommand struct {
	Email     string `json:"email" validate:"required|email"`
	Password  string `json:"password" validate:"required|min=8"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}
