// Package model
package model

import vo "github.com/paladignus/actajus/internal/shared/domain/value_object"

type LoginNormalized struct {
	Email     vo.Email
	Password  string
	IP        string
	UserAgent string
}
