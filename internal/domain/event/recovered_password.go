// Package event
package event

import "log"

// PasswordResetRequestedEvent é o evento disparado quando um usuário
// solicita a redefinição de senha.
//
// Este evento carrega todas as informações necessárias para que os
// handlers possam executar suas ações (ex: enviar email).
type PasswordResetRequestedEvent struct {
	// BaseEvent contém campos comuns: Name, Timestamp, Version, etc
	BaseEvent

	// Email do usuário que solicitou a redefinição
	UserEmail string

	// Token gerado para redefinição de senha
	// Será usado na URL enviada por email
	ResetToken string

	// URL completa para redefinição
	// Exemplo: "https://seusite.com/reset-password?token=abc123"
	ResetURL string

	// Nome do usuário (opcional, para personalização do email)
	UserName string
}

// NewPasswordResetRequestedEvent cria um novo evento de solicitação
// de redefinição de senha.
//
// Parâmetros:
//   - userID: ID do usuário no sistema
//   - email: email do usuário
//   - resetToken: token de redefinição gerado
//   - resetURL: URL completa para redefinição
//   - userName: nome do usuário (pode ser vazio)
//
// Retorna:
//   - PasswordResetRequestedEvent preenchido e pronto para publicação
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

// init registra automaticamente todos os eventos deste arquivo no registry global.
// Isso garante que os eventos podem ser deserializados corretamente.
func init() {
	// Registra PasswordResetRequestedEvent
	Register("user.password_reset_requested", func() Event {
		return &PasswordResetRequestedEvent{}
	})

	log.Println("[EventRegistry] Registered: user.password_reset_requested")

	// Quando criar novos eventos, adicione aqui:
	// Register("user.created", func() Event {
	//     return &UserCreatedEvent{}
	// })
}

// Exemplo de outro evento que você pode criar no futuro:
//
// type UserCreatedEvent struct {
//     BaseEvent
//     UserEmail string
//     UserName  string
// }
//
// func NewUserCreatedEvent(userID, email, name string) UserCreatedEvent {
//     return UserCreatedEvent{
//         BaseEvent: NewBaseEvent("user.created", userID, "v1"),
//         UserEmail: email,
//         UserName:  name,
//     }
// }

// type PasswordResetRequestedEvent struct {
// 	BaseEvent
// 	UserEmail  string
// 	ResetToken string
// 	ResetURL   string
// 	UserName   string
// }
//
// func NewPasswordResetRequestedEvent(
// 	userID string,
// 	email string,
// 	resetToken string,
// 	resetURL string,
// 	userName string,
// ) PasswordResetRequestedEvent {
// 	return PasswordResetRequestedEvent{
// 		// Cria o BaseEvent com nome e versão do evento
// 		BaseEvent: NewBaseEvent(
// 			"user.password_reset_requested", // Nome do evento
// 			userID,                          // ID da entidade (usuário)
// 			"v1",                            // Versão do schema
// 		),
// 		UserEmail:  email,
// 		ResetToken: resetToken,
// 		ResetURL:   resetURL,
// 		UserName:   userName,
// 	}
// }
//
// func (r PasswordResetRequestedEvent) Serialize() ([]byte, error) {
// 	return json.Marshal(r)
// }
//
// func (r PasswordResetRequestedEvent) Deserialize(data []byte) (PasswordResetRequestedEvent, error) {
// 	if err := json.Unmarshal(data, &r); err != nil {
// 		return PasswordResetRequestedEvent{}, err
// 	}
// 	return r, nil
// }
