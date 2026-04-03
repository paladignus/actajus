// Package adapter centralizes presentation-layer mappings.
package adapter

import (
	"github.com/paladignus/actajus/internal/shared/domain"
)

type ValidationViolationPayload struct {
	Path string            `json:"path"`
	Code string            `json:"code"`
	Meta map[string]string `json:"meta,omitempty"`
}

type ValidationErrorPayload struct {
	Violations []ValidationViolationPayload `json:"violations"`
}

func ValidationPayload(err domain.ValidationError) ValidationErrorPayload {
	payload := ValidationErrorPayload{
		Violations: make([]ValidationViolationPayload, 0, len(err.Violations)),
	}
	for _, violation := range err.Violations {
		meta := violation.Meta
		if meta == nil {
			meta = map[string]string{}
		}
		payload.Violations = append(payload.Violations, ValidationViolationPayload{
			Path: violation.Path,
			Code: string(violation.Code),
			Meta: meta,
		})
	}
	return payload
}

func ValidationErrorsByPath(err domain.ValidationError) map[string][]ValidationViolationPayload {
	grouped := make(map[string][]ValidationViolationPayload, len(err.Violations))
	for _, violation := range ValidationPayload(err).Violations {
		grouped[violation.Path] = append(grouped[violation.Path], violation)
	}
	return grouped
}
