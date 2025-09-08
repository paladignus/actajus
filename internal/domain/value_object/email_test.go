package valueobject

import (
	"testing"
)

type Email string

func (e Email) Value() string {
	return string(e)
}

func (e Email) IsEmpty() bool {
	return len(e) == 0
}

func (e Email) IsValid() bool {
	return !e.IsEmpty()
}

func TestEmail(t *testing.T) {
	t.Run("empty email should be invalid", func(t *testing.T) {
		sut := Email("")
		if sut.IsValid() {
			t.Errorf("expected empty email to be invalid")
		}
	})
}
