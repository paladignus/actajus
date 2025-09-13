// Package valueobject
package valueobject

import (
	"regexp"
	"testing"
)

type Text string

func (t Text) Value() string {
	return string(t)
}

func (t Text) IsValid() bool {
	re := regexp.MustCompile("^[A-Za-záàâãéèêíïóôõöúçñÁÀÂÃÉÈÍÏÓÔÕÖÚÇÑ ]+$")
	return re.MatchString(t.Value())
}

func TestText(t *testing.T) {
	sut := Text("M4rc3l0@#$%¨&*()")
	t.Run("should return error if text contains numbers and special characters", func(t *testing.T) {
		if sut.IsValid() {
			t.Errorf("Expected error for text with numbers and special characters, but got nil")
		}
	})
	sut = Text("Marcelo Bento Pereira")
	t.Run("should return nil if text is valid", func(t *testing.T) {
		if !sut.IsValid() {
			t.Errorf("Expected nil for valid text, but got error")
		}
	})
}
