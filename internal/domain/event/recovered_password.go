// Package event
package event

type RecoveredPassword struct {
	Email string
}

func (r RecoveredPassword) Subject() string {
	return "user.recovered_password"
}

func (r RecoveredPassword) Serialize() ([]byte, error) {
	return []byte(r.Email), nil
}
