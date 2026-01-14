// Package entity
package entity

import "github.com/paladignus/actajus/internal/application/dto"

type Phone struct {
	IDPhone    int
	Number     string
	Kind       string
	Department string
}

func NewPhone(phone dto.Phone) Phone {
	return Phone{
		Number:     phone.Number,
		Kind:       phone.Kind,
		Department: phone.Department,
	}
}
