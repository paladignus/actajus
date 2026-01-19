// Package valueobject
package valueobject

import "regexp"

type ZIP string

func (z ZIP) Value() string {
	return string(z)
}

func (z ZIP) IsValid() bool {
	regex := regexp.MustCompile(`^[0-9]{5}-[0-9]{3}$`)
	return regex.MatchString(z.Value())
}
