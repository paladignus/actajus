// Package event
package event

import "encoding/json"

type PasswordResetRequestedEvent struct {
	BaseEvent
	UserEmail  string
	ResetToken string
	ResetURL   string
	UserName   string
}

func NewPasswordResetRequestedEvent(
	userID string,
	email string,
	resetToken string,
	resetURL string,
	userName string,
) PasswordResetRequestedEvent {
	return PasswordResetRequestedEvent{
		// Cria o BaseEvent com nome e versão do evento
		BaseEvent: NewBaseEvent(
			"user.password_reset_requested", // Nome do evento
			userID,                          // ID da entidade (usuário)
			"v1",                            // Versão do schema
		),
		UserEmail:  email,
		ResetToken: resetToken,
		ResetURL:   resetURL,
		UserName:   userName,
	}
}

func (r PasswordResetRequestedEvent) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

func (r PasswordResetRequestedEvent) Deserialize(data []byte) (PasswordResetRequestedEvent, error) {
	if err := json.Unmarshal(data, &r); err != nil {
		return PasswordResetRequestedEvent{}, err
	}
	return r, nil
}
