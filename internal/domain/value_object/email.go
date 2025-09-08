package valueobject

import "regexp"

type Email string

func (e Email) Value() string {
	return string(e)
}

func (e Email) IsEmpty() bool {
	return len(e) == 0
}

func (e Email) Equals(other Email) bool {
	return e.Value() == other.Value()
}

func (e Email) IsValid() bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return !e.IsEmpty() && re.MatchString(e.Value())
}
