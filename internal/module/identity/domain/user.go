// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

// User representa uma entidade de domínio do contexto de Identity.
// Como Entidade Rica, ela contém lógica de negócio e protege seus invariantes.
type User struct {
	id           vo.ID
	primaryEmail vo.Email
	passwordHash vo.PasswordHash
	isBlocked    bool
	createdAt    time.Time
	updatedAt    time.Time
}

// UserBuilder segue o padrão Builder para construção segura de User.
type UserBuilder struct {
	u *User
}

// NewUserBuilder cria um novo builder para User.
func NewUserBuilder() *UserBuilder {
	now := time.Now()
	return &UserBuilder{
		u: &User{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (b *UserBuilder) WithID(id int64) *UserBuilder {
	b.u.id = vo.ID(id)
	return b
}

func (b *UserBuilder) WithPrimaryEmail(email string) *UserBuilder {
	b.u.primaryEmail = vo.Email(email)
	return b
}

func (b *UserBuilder) WithPasswordHash(hash string) *UserBuilder {
	b.u.passwordHash = vo.PasswordHash(hash)
	return b
}

func (b *UserBuilder) WithIsBlocked(blocked bool) *UserBuilder {
	b.u.isBlocked = blocked
	return b
}

func (b *UserBuilder) WithCreatedAt(createdAt time.Time) *UserBuilder {
	b.u.createdAt = createdAt
	return b
}

func (b *UserBuilder) WithUpdatedAt(updatedAt time.Time) *UserBuilder {
	b.u.updatedAt = updatedAt
	return b
}

// Build valida e constrói a entidade User.
// Retorna erro se os invariantes de domínio forem violados.
func (b *UserBuilder) Build() (*User, error) {
	if err := b.u.validate(); err != nil {
		return nil, err
	}
	return b.u, nil
}

// validate verifica os invariantes de domínio da entidade.
func (u *User) validate() error {
	if u.primaryEmail.IsEmpty() || !u.primaryEmail.IsValid() {
		return domain.NewFieldError("email", "email is invalid")
	}
	return nil
}

// ID retorna o identificador único do usuário.
func (u *User) ID() vo.ID { return u.id }

// PrimaryEmail retorna o email principal do usuário.
func (u *User) PrimaryEmail() vo.Email { return u.primaryEmail }

// PasswordHash retorna o hash da senha (apenas para leitura).
func (u *User) PasswordHash() vo.PasswordHash { return u.passwordHash }

// IsBlocked retorna se o usuário está bloqueado.
func (u *User) IsBlocked() bool { return u.isBlocked }

// CreatedAt retorna a data de criação do usuário.
func (u *User) CreatedAt() time.Time { return u.createdAt }

// UpdatedAt retorna a data da última atualização do usuário.
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// SetID define o ID do usuário apenas uma vez.
// Protege o invariante: ID não pode ser alterado após definição.
func (u *User) SetID(id vo.ID) error {
	if u.id != 0 {
		return domain.NewFieldError("id", "user id is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "user id is invalid")
	}
	u.id = id
	return nil
}

// MatchesPassword verifica se a senha fornecida corresponde ao hash armazenado.
// Este método encapsula a lógica de comparação de senha no domínio.
func (u *User) MatchesPassword(password string, hasher PasswordHasher) bool {
	return hasher.Compare(u.passwordHash.Value(), password) == nil
}

// ChangePassword altera a senha do usuário após validar as regras de negócio.
// Regras aplicadas:
//   - Usuário não pode estar bloqueado
//   - Nova senha não pode ser vazia
//   - Atualiza o timestamp de updatedAt
//
// Este é um exemplo de comportamento de domínio: a entidade protege
// seu estado e aplica regras antes de permitir a mudança.
func (u *User) ChangePassword(newPassword string, hasher PasswordHasher) error {
	if u.isBlocked {
		return ErrUserBlocked
	}

	// Valida que a nova senha não é vazia
	if newPassword == "" {
		return domain.NewFieldError("password", "password cannot be empty")
	}

	// Gera o novo hash
	hashed, err := hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	u.passwordHash = vo.PasswordHash(hashed)
	u.updatedAt = time.Now()
	return nil
}

// Block bloqueia o usuário.
// Uma vez bloqueado, o usuário não pode realizar operações sensíveis.
// Esta operação é idempotente - bloquear um usuário já bloqueado não causa erro.
func (u *User) Block() {
	if u.isBlocked {
		return
	}
	u.isBlocked = true
	u.updatedAt = time.Now()
}

// Unblock desbloqueia o usuário.
// Esta operação é idempotente - desbloquear um usuário já desbloqueado não causa erro.
func (u *User) Unblock() {
	if !u.isBlocked {
		return
	}
	u.isBlocked = false
	u.updatedAt = time.Now()
}

// EmailEquals verifica se o email fornecido é igual ao email principal do usuário.
// Útil para validações de login e recuperação de senha.
func (u *User) EmailEquals(email string) bool {
	return u.primaryEmail.Equals(vo.Email(email))
}
