// Package mapper
package mapper

import "github.com/paladignus/actajus/internal/shared/presentation/validation"

func shouldAddEmailViolation(vs []validation.Violation) bool {
	for _, v := range vs {
		if v.Path == "email" {
			return false
		}
	}
	return true
}

// func hasViolation(vs []validation.Violation, path string) bool {
// 	for _, v := range vs {
// 		if v.Path == path {
// 			return true
// 		}
// 	}
// 	return false
// }
