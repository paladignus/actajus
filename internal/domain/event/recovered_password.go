// Package event
package event

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
