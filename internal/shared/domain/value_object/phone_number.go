// Package valueobject
package valueobject

import "regexp"

type PhoneNumber string

func (p PhoneNumber) Value() string {
	return string(p)
}

func (p PhoneNumber) IsValid() bool {
	re := regexp.MustCompile(`^\+55[0-9]{10,11}$`)
	return re.MatchString(p.Value())
}
