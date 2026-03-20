// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

// Session representa uma entidade de domínio do contexto de Identity.
// Uma sessão é criada quando um usuário faz login e contém tokens de refresh.
// Como Entidade Rica, ela protege invariantes como:
//   - Não pode pertencer a um usuário inválido
//   - Não pode ser usada após expiração ou revogação
//   - Rotação de tokens deve ser atômica para prevenir replay attacks
type Session struct {
	id          vo.ID
	idUser      vo.ID
	refreshHash [32]byte
	expiresAt   time.Time
	revokedAt   *time.Time
	rotatedAt   *time.Time
	ip          string
	userAgent   string
	createdAt   time.Time
	updatedAt   time.Time
}

// SessionBuilder segue o padrão Builder para construção segura de Session.
type SessionBuilder struct {
	s *Session
}

// NewSessionBuilder cria um novo builder para Session.
func NewSessionBuilder() *SessionBuilder {
	now := time.Now()
	return &SessionBuilder{
		&Session{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (b *SessionBuilder) WithID(id int64) *SessionBuilder {
	b.s.id = vo.ID(id)
	return b
}

func (b *SessionBuilder) WithIDUser(id int64) *SessionBuilder {
	b.s.idUser = vo.ID(id)
	return b
}

func (b *SessionBuilder) WithRefreshHash(h [32]byte) *SessionBuilder {
	b.s.refreshHash = h
	return b
}

func (b *SessionBuilder) WithExpiresAt(t time.Time) *SessionBuilder {
	b.s.expiresAt = t
	return b
}

func (b *SessionBuilder) WithRevokedAt(t *time.Time) *SessionBuilder {
	b.s.revokedAt = t
	return b
}

func (b *SessionBuilder) WithRotatedAt(t *time.Time) *SessionBuilder {
	b.s.rotatedAt = t
	return b
}

func (b *SessionBuilder) WithIP(ip string) *SessionBuilder {
	b.s.ip = ip
	return b
}

func (b *SessionBuilder) WithUserAgent(ua string) *SessionBuilder {
	b.s.userAgent = ua
	return b
}

func (b *SessionBuilder) WithCreatedAt(t time.Time) *SessionBuilder {
	b.s.createdAt = t
	return b
}

func (b *SessionBuilder) WithUpdatedAt(t time.Time) *SessionBuilder {
	b.s.updatedAt = t
	return b
}

// Build valida e constrói a entidade Session.
// Retorna erro se os invariantes de domínio forem violados.
func (b *SessionBuilder) Build() (*Session, error) {
	if err := b.s.validate(); err != nil {
		return nil, err
	}
	return b.s, nil
}

// validate verifica os invariantes de domínio da entidade.
func (s *Session) validate() error {
	if s.idUser == 0 {
		return domain.NewFieldError("user_id", "user id is invalid")
	}
	if s.expiresAt.IsZero() {
		return domain.NewFieldError("expires_at", "expires_at is required")
	}
	return nil
}

// ID retorna o identificador único da sessão.
func (s *Session) ID() vo.ID { return s.id }

// IDUser retorna o identificador do usuário dono desta sessão.
func (s *Session) IDUser() vo.ID { return s.idUser }

// RefreshHash retorna o hash do token de refresh (apenas para leitura).
func (s *Session) RefreshHash() [32]byte { return s.refreshHash }

// ExpiresAt retorna a data de expiração da sessão.
func (s *Session) ExpiresAt() time.Time { return s.expiresAt }

// RevokedAt retorna a data de revogação (nil se não revogada).
func (s *Session) RevokedAt() *time.Time { return s.revokedAt }

// RotatedAt retorna a data da última rotação de token (nil se nunca rotacionada).
func (s *Session) RotatedAt() *time.Time { return s.rotatedAt }

// IP retorna o endereço IP de origem da sessão.
func (s *Session) IP() string { return s.ip }

// UserAgent retorna o User-Agent de origem da sessão.
func (s *Session) UserAgent() string { return s.userAgent }

// CreatedAt retorna a data de criação da sessão.
func (s *Session) CreatedAt() time.Time { return s.createdAt }

// UpdatedAt retorna a data da última atualização da sessão.
func (s *Session) UpdatedAt() time.Time { return s.updatedAt }

// IsRevoked verifica se a sessão foi revogada.
func (s *Session) IsRevoked() bool { return s.revokedAt != nil }

// IsExpired verifica se a sessão expirou em relação ao tempo fornecido.
func (s *Session) IsExpired(now time.Time) bool { return !now.Before(s.expiresAt) }

// IsActive verifica se a sessão está ativa (não revogada e não expirada).
// Este método encapsula a lógica de verificação de estado da sessão.
func (s *Session) IsActive(now time.Time) bool {
	return !s.IsRevoked() && !s.IsExpired(now)
}

// SetID define o ID da sessão apenas uma vez.
// Protege o invariante: ID não pode ser alterado após definição.
func (s *Session) SetID(id int64) error {
	if s.id != 0 {
		return domain.NewFieldError("id", "session id is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "session id is invalid")
	}
	s.id = vo.ID(id)
	return nil
}

// Revoke revoga a sessão, impedindo seu uso futuro.
// Esta operação é idempotente - revogar uma sessão já revogada retorna erro.
// A revogação atualiza o timestamp updatedAt para rastrear quando ocorreu.
func (s *Session) Revoke(now time.Time) error {
	if s.id == 0 {
		return domain.NewFieldError("id", "session id is invalid")
	}
	if s.IsRevoked() {
		return domain.NewFieldError("revoked_at", "session already revoked")
	}
	s.revokedAt = &now
	s.updatedAt = now
	return nil
}

// Rotate atualiza o token de refresh da sessão de forma atômica.
// A rotação de tokens é uma medida de segurança para prevenir replay attacks.
// Regras aplicadas:
//   - Sessão não pode estar revogada
//   - Sessão deve ter ID válido
//   - Atualiza hash, expiração e timestamp de rotação
//
// Este método é chamado durante o refresh de tokens para garantir
// que cada token de refresh só possa ser usado uma vez.
func (s *Session) Rotate(hash [32]byte, exp time.Time, now time.Time) error {
	if s.id == 0 {
		return domain.NewFieldError("id", "session id is invalid")
	}
	if s.IsRevoked() {
		return ErrSessionRevoked
	}
	s.refreshHash = hash
	s.expiresAt = exp
	s.rotatedAt = &now
	s.updatedAt = now
	return nil
}

// CanRotate verifica se a sessão pode ter seu token rotacionado.
// Útil para validações prévias antes de chamar Rotate.
func (s *Session) CanRotate(now time.Time) bool {
	return s.IsActive(now)
}

// TimeToLive retorna o tempo restante até a expiração da sessão.
// Retorna 0 se a sessão já expirou ou foi revogada.
func (s *Session) TimeToLive(now time.Time) time.Duration {
	if !s.IsActive(now) {
		return 0
	}
	return s.expiresAt.Sub(now)
}

// IsAboutToExpire verifica se a sessão está prestes a expirar.
// threshold: tempo mínimo aceitável antes da expiração (ex: 5 minutos).
func (s *Session) IsAboutToExpire(now time.Time, threshold time.Duration) bool {
	if !s.IsActive(now) {
		return false
	}
	return s.TimeToLive(now) < threshold
}
