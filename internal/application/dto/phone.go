// Package dto
package dto

type Phone struct {
	IDPhone    int    `json:"id_phone"`
	Number     string `json:"number"`
	Kind       string `json:"kind"`
	Department string `json:"department"`
}
