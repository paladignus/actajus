// Package event
package event

type PasswordResetRequestedEvent struct {
	Event
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
		Event: NewEvent(
			"user.password_reset_requested",
			userID,
			"v1",
		),
		UserEmail:  email,
		ResetToken: resetToken,
		ResetURL:   resetURL,
		UserName:   userName,
	}
}
