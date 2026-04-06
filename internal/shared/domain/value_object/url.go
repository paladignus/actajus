// Package valueobject
package valueobject

import "regexp"

type URL string

func (f URL) Value() string {
	return string(f)
}

func (f URL) IsValid() bool {
	re := regexp.MustCompile(`^(http|https)://(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)
	return re.MatchString(f.Value())
}
