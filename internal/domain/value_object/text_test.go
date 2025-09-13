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
	re := regexp.MustCompile(`[\W_]|[0-9]`)
	return !re.MatchString(t.Value())
}

func TestText(t *testing.T) {
	sut := Text("M4rc3l0@#$%¨&*()")
	t.Run("should return error if text contains numbers and special characters", func(t *testing.T) {
		if sut.IsValid() {
			t.Errorf("Expected error for text with numbers and special characters, but got nil")
		}
	})
}
