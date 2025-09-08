package valueobject

import (
	"regexp"
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
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return !e.IsEmpty() && re.MatchString(e.Value())
}

func TestEmail(t *testing.T) {
	t.Run("empty email should be invalid", func(t *testing.T) {
		sut := Email("")
		if sut.IsValid() {
			t.Errorf("expected empty email to be invalid")
		}
	})
	t.Run("invalid email should be invalid", func(t *testing.T) {
		sut := Email("invalid-email")
		if sut.IsValid() {
			t.Errorf("expected invalid email to be invalid")
		}
	})
}
