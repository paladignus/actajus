// Package dto
package dto

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	IP        string `json:"-"`
	UserAgent string `json:"-"`
}
