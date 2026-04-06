// Package adapter
package adapter

import (
	"github.com/paladignus/actajus/internal/shared/domain"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
)

func ViolationsToDomainError(vs []validation.Violation) error {
	if len(vs) == 0 {
		return nil
	}
	violations := make([]domain.Violation, 0, len(vs))
	for _, v := range vs {
		violations = append(violations, domain.Violation{
			Path: v.Path,
			Code: domain.ViolationCode(v.Code),
			Meta: v.Meta,
		})
	}
	return domain.NewValidationError(violations)
}
