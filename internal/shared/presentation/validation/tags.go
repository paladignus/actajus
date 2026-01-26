// Package validation
package validation

import (
	"strings"
)

func parseRule(rule string) (name string, param int) {
	parts := strings.Split(rule, "=")
	name = parts[0]
	if len(parts) == 2 {
		param = atoi(parts[1])
	}
	return
}
