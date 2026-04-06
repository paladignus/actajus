// Package valueobject
package valueobject

import "regexp"

type Text string

func (t Text) Value() string {
	return string(t)
}

func (t Text) IsValid() bool {
	re := regexp.MustCompile("^[A-Za-záàâãéèêíïóôõöúçñÁÀÂÃÉÈÍÏÓÔÕÖÚÇÑ ]+$")
	return re.MatchString(t.Value())
}
