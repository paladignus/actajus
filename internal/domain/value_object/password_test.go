package valueobject

import (
	"testing"
)

type Password string

func (p Password) Value() string {
	return string(p)
}

func (p Password) IsEmpty() bool {
	return len(p) == 0
}

func (p Password) IsValid() bool {
	return !p.IsEmpty()
}

func TestPassword(t *testing.T) {
	sut := Password("")
	t.Run("empty password should be invalid", func(t *testing.T) {
		if sut.IsValid() {
			t.Errorf("expected empty password to be invalid")
		}
	})
}
