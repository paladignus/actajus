// Package valueobject
package valueobject

import "regexp"

type File string

func (f File) Value() string {
	return string(f)
}

func (f File) IsValid() bool {
	re := regexp.MustCompile(`^[\w,\s-]+\.(png|jpg|gif|jpeg|pdf)$`)
	return re.MatchString(f.Value())
}
