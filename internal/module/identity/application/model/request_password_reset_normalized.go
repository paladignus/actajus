// Package model
package model

import vo "github.com/paladignus/actajus/internal/shared/domain/value_object"

type RequestPasswordResetNormalized struct {
	Email vo.Email
}
