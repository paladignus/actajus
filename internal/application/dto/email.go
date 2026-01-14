// Package dto
package dto

type SendEmailRecoverPasswordInput struct {
	Address string
	URL     string
}

type Email struct {
	IDEmails uint   `json:"idemails"`
	Address  string `json:"address"`
}
