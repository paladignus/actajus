// Package valueobject
package valueobject

import "regexp"

type Date string

func (d Date) Value() string {
	return string(d)
}

func (d Date) IsValid() bool {
	re := regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/\d{4}$`)
	return re.MatchString(d.Value())
}
