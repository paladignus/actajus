// Package dto
package dto

import "time"

// ============================================================================
// Identity Events
// ============================================================================

// UserCreated é disparado quando um usuário é criado
type UserCreated struct {
	IDUser    int64  `json:"id_user"`
	Email     string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (e UserCreated) EventType() string   { return "user.created" }
func (e UserCreated) OccurredAt() time.Time { return e.CreatedAt }
func (e UserCreated) AggregateID() string  { return string(rune(e.IDUser)) }
func (e UserCreated) Version() string     { return "1.0" }

// UserPasswordChanged é disparado quando a senha de um usuário é alterada
type UserPasswordChanged struct {
	IDUser    int64     `json:"id_user"`
	ChangedAt time.Time `json:"changed_at"`
}

func (e UserPasswordChanged) EventType() string   { return "user.password_changed" }
func (e UserPasswordChanged) OccurredAt() time.Time { return e.ChangedAt }
func (e UserPasswordChanged) AggregateID() string  { return string(rune(e.IDUser)) }
func (e UserPasswordChanged) Version() string     { return "1.0" }

// UserBlocked é disparado quando um usuário é bloqueado
type UserBlocked struct {
	IDUser    int64     `json:"id_user"`
	BlockedAt time.Time `json:"blocked_at"`
}

func (e UserBlocked) EventType() string   { return "user.blocked" }
func (e UserBlocked) OccurredAt() time.Time { return e.BlockedAt }
func (e UserBlocked) AggregateID() string  { return string(rune(e.IDUser)) }
func (e UserBlocked) Version() string     { return "1.0" }

// SessionCreated é disparado quando uma sessão é criada (login)
type SessionCreated struct {
	IDSession int64     `json:"id_session"`
	IDUser    int64     `json:"id_user"`
	CreatedAt time.Time `json:"created_at"`
}

func (e SessionCreated) EventType() string   { return "session.created" }
func (e SessionCreated) OccurredAt() time.Time { return e.CreatedAt }
func (e SessionCreated) AggregateID() string  { return string(rune(e.IDSession)) }
func (e SessionCreated) Version() string     { return "1.0" }

// SessionRevoked é disparado quando uma sessão é revogada (logout)
type SessionRevoked struct {
	IDSession int64     `json:"id_session"`
	RevokedAt time.Time `json:"revoked_at"`
}

func (e SessionRevoked) EventType() string   { return "session.revoked" }
func (e SessionRevoked) OccurredAt() time.Time { return e.RevokedAt }
func (e SessionRevoked) AggregateID() string  { return string(rune(e.IDSession)) }
func (e SessionRevoked) Version() string     { return "1.0" }

// ============================================================================
// Company Events
// ============================================================================

// CompanyCreated é disparado quando uma empresa é criada
type CompanyCreated struct {
	IDCompany   int64     `json:"id_company"`
	Name        string    `json:"name"`
	CNPJ        string    `json:"cnpj"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

func (e CompanyCreated) EventType() string   { return "company.created" }
func (e CompanyCreated) OccurredAt() time.Time { return e.CreatedAt }
func (e CompanyCreated) AggregateID() string  { return string(rune(e.IDCompany)) }
func (e CompanyCreated) Version() string     { return "1.0" }

// CompanyUpdated é disparado quando uma empresa é atualizada
type CompanyUpdated struct {
	IDCompany int64     `json:"id_company"`
	UpdatedAt time.Time `json:"updated_at"`
	Changes   []string  `json:"changes,omitempty"`
}

func (e CompanyUpdated) EventType() string   { return "company.updated" }
func (e CompanyUpdated) OccurredAt() time.Time { return e.UpdatedAt }
func (e CompanyUpdated) AggregateID() string  { return string(rune(e.IDCompany)) }
func (e CompanyUpdated) Version() string     { return "1.0" }

// CompanyDeleted é disparado quando uma empresa é excluída
type CompanyDeleted struct {
	IDCompany int64     `json:"id_company"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (e CompanyDeleted) EventType() string   { return "company.deleted" }
func (e CompanyDeleted) OccurredAt() time.Time { return e.DeletedAt }
func (e CompanyDeleted) AggregateID() string  { return string(rune(e.IDCompany)) }
func (e CompanyDeleted) Version() string     { return "1.0" }

// CompanyAddressAdded é disparado quando um endereço é associado a uma empresa
type CompanyAddressAdded struct {
	IDCompany int64     `json:"id_company"`
	IDAddress int64     `json:"id_address"`
	AddedAt   time.Time `json:"added_at"`
}

func (e CompanyAddressAdded) EventType() string   { return "company.address_added" }
func (e CompanyAddressAdded) OccurredAt() time.Time { return e.AddedAt }
func (e CompanyAddressAdded) AggregateID() string  { return string(rune(e.IDCompany)) }
func (e CompanyAddressAdded) Version() string     { return "1.0" }

// CompanyPhoneAdded é disparado quando um telefone é associado a uma empresa
type CompanyPhoneAdded struct {
	IDCompany int64     `json:"id_company"`
	IDPhone   int64     `json:"id_phone"`
	AddedAt   time.Time `json:"added_at"`
}

func (e CompanyPhoneAdded) EventType() string   { return "company.phone_added" }
func (e CompanyPhoneAdded) OccurredAt() time.Time { return e.AddedAt }
func (e CompanyPhoneAdded) AggregateID() string  { return string(rune(e.IDCompany)) }
func (e CompanyPhoneAdded) Version() string     { return "1.0" }

// CompanyEmailAdded é disparado quando um email é associado a uma empresa
type CompanyEmailAdded struct {
	IDCompany int64     `json:"id_company"`
	IDEmail   int64     `json:"id_email"`
	AddedAt   time.Time `json:"added_at"`
}

func (e CompanyEmailAdded) EventType() string   { return "company.email_added" }
func (e CompanyEmailAdded) OccurredAt() time.Time { return e.AddedAt }
func (e CompanyEmailAdded) AggregateID() string  { return string(rune(e.IDCompany)) }
func (e CompanyEmailAdded) Version() string     { return "1.0" }
