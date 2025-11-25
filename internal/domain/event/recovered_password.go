// Package event
package event

type PasswordResetRequestedEvent struct {
	Event
	UserEmail string
	ResetURL  string
}

func NewPasswordResetRequestedEvent(
	userID string,
	email string,
	resetURL string,
) PasswordResetRequestedEvent {
	return PasswordResetRequestedEvent{
		Event: NewEvent(
			"user.password_reset_requested",
			userID,
			"v1",
		),
		UserEmail: email,
		ResetURL:  resetURL,
	}
}
