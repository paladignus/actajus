// Package dto
package dto

type CreatePhoneRequest struct {
	Number     string `json:"number"`
	Kind       string `json:"kind"`
	Department string `json:"department"`
}
