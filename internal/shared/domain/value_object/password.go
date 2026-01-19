// Package valueobject
package valueobject

import "regexp"

type Password string

func (p Password) Value() string {
	return string(p)
}

func (p Password) IsEmpty() bool {
	return len(p) == 0
}

func (p Password) IsValid() bool {
	number := regexp.MustCompile(`[0-9]`)
	lower := regexp.MustCompile(`[a-z]`)
	upper := regexp.MustCompile(`[A-Z]`)
	special := regexp.MustCompile(`[\W_]`)
	return !p.IsEmpty() &&
		len(p.Value()) > 7 &&
		number.MatchString(p.Value()) &&
		lower.MatchString(p.Value()) &&
		upper.MatchString(p.Value()) &&
		special.MatchString(p.Value())
}
