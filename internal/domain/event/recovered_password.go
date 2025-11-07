// Package event
package event

import "encoding/json"

type RecoveredPassword struct {
	Email string
	URL   string
}

func (r RecoveredPassword) Subject() string {
	return "user.created"
}

func (r RecoveredPassword) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

func (r RecoveredPassword) Deserialize(data []byte) (RecoveredPassword, error) {
	if err := json.Unmarshal(data, &r); err != nil {
		return RecoveredPassword{}, err
	}
	return r, nil
}
