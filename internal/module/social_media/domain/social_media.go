// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

// SocialMedia representa uma entidade de domínio compartilhada entre Company e Person.
// Como Entidade Rica, ela contém lógica de negócio e protege seus invariantes.
type SocialMedia struct {
	id        vo.ID
	idCompany vo.ID
	platform  vo.Text
	url       vo.URL
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

// SocialMediaBuilder segue o padrão Builder para construção segura de SocialMedia.
type SocialMediaBuilder struct {
	socialMedia *SocialMedia
}

// NewSocialMediaBuilder cria um novo builder para SocialMedia.
func NewSocialMediaBuilder() *SocialMediaBuilder {
	now := time.Now()
	return &SocialMediaBuilder{
		socialMedia: &SocialMedia{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (s *SocialMediaBuilder) WithID(id int64) *SocialMediaBuilder {
	s.socialMedia.id = vo.ID(id)
	return s
}

func (s *SocialMediaBuilder) WithIDCompany(id int64) *SocialMediaBuilder {
	s.socialMedia.idCompany = vo.ID(id)
	return s
}

func (s *SocialMediaBuilder) WithPlatform(platform string) *SocialMediaBuilder {
	s.socialMedia.platform = vo.Text(platform)
	return s
}

func (s *SocialMediaBuilder) WithURL(url string) *SocialMediaBuilder {
	s.socialMedia.url = vo.URL(url)
	return s
}

func (s *SocialMediaBuilder) WithCreatedAt(createdAt time.Time) *SocialMediaBuilder {
	s.socialMedia.createdAt = createdAt
	return s
}

func (s *SocialMediaBuilder) WithUpdatedAt(updatedAt time.Time) *SocialMediaBuilder {
	s.socialMedia.updatedAt = updatedAt
	return s
}

func (s *SocialMediaBuilder) WithDeletedAt(deletedAt *time.Time) *SocialMediaBuilder {
	s.socialMedia.deletedAt = deletedAt
	return s
}

// Build valida e constrói a entidade SocialMedia.
func (s *SocialMediaBuilder) Build() (*SocialMedia, error) {
	if err := s.socialMedia.validate(); err != nil {
		return nil, err
	}
	return s.socialMedia, nil
}

// validate verifica os invariantes de domínio da entidade.
func (s *SocialMedia) validate() error {
	if !s.platform.IsValid() {
		return domain.NewFieldError("platform", "platform is invalid")
	}
	if !s.url.IsValid() {
		return domain.NewFieldError("url", "url is invalid")
	}
	return nil
}

// ID retorna o identificador único da mídia social.
func (s *SocialMedia) ID() vo.ID { return s.id }

// IDCompany retorna o identificador da empresa associada.
func (s *SocialMedia) IDCompany() vo.ID { return s.idCompany }

// Platform retorna a plataforma da mídia social (ex: "LinkedIn", "Twitter").
func (s *SocialMedia) Platform() vo.Text { return s.platform }

// URL retorna a URL do perfil da mídia social.
func (s *SocialMedia) URL() vo.URL { return s.url }

// CreatedAt retorna a data de criação da mídia social.
func (s *SocialMedia) CreatedAt() time.Time { return s.createdAt }

// UpdatedAt retorna a data da última atualização da mídia social.
func (s *SocialMedia) UpdatedAt() time.Time { return s.updatedAt }

// DeletedAt retorna a data de exclusão (nil se não excluída).
func (s *SocialMedia) DeletedAt() *time.Time { return s.deletedAt }

// SetID define o ID da mídia social apenas uma vez.
func (s *SocialMedia) SetID(id int64) error {
	if s.id != 0 {
		return domain.NewFieldError("id", "social media ID is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "social media ID is invalid")
	}
	s.id = vo.ID(id)
	return nil
}

// SetCompanyID define o ID da empresa associada.
func (s *SocialMedia) SetCompanyID(id int64) error {
	if s.idCompany != 0 {
		return domain.NewFieldError("id_company", "company ID is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id_company", "company ID is invalid")
	}
	s.idCompany = vo.ID(id)
	return nil
}

// IsDeleted verifica se a mídia social foi excluída logicamente.
func (s *SocialMedia) IsDeleted() bool {
	return s.deletedAt != nil
}

// Delete exclui logicamente a mídia social.
func (s *SocialMedia) Delete() error {
	if s.id == 0 {
		return domain.NewFieldError("id", "social media ID is invalid")
	}
	if s.IsDeleted() {
		return domain.NewFieldError("deleted_at", "social media is already deleted")
	}
	now := time.Now()
	s.deletedAt = &now
	s.updatedAt = now
	return nil
}

// Restore restaura uma mídia social excluída logicamente.
func (s *SocialMedia) Restore() {
	if !s.IsDeleted() {
		return
	}
	s.deletedAt = nil
	s.updatedAt = time.Now()
}

// UpdatePlatform atualiza a plataforma da mídia social.
func (s *SocialMedia) UpdatePlatform(platform string) error {
	if !vo.Text(platform).IsValid() {
		return domain.NewFieldError("platform", "platform is invalid")
	}
	s.platform = vo.Text(platform)
	s.updatedAt = time.Now()
	return nil
}

// UpdateURL atualiza a URL do perfil da mídia social.
func (s *SocialMedia) UpdateURL(url string) error {
	if !vo.URL(url).IsValid() {
		return domain.NewFieldError("url", "url is invalid")
	}
	s.url = vo.URL(url)
	s.updatedAt = time.Now()
	return nil
}
