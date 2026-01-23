// Package domain
package domain

import (
	"time"

	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type SocialMedia struct {
	id        uint
	idCompany uint
	name      vo.Text
	url       vo.URL
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

type SocialMediaBuilder struct {
	socialMedia *SocialMedia
	errors      []error
}

func NewSocialMediaBuilder() *SocialMediaBuilder {
	now := time.Now()
	return &SocialMediaBuilder{
		socialMedia: &SocialMedia{
			createdAt: now,
			updatedAt: now,
		},
		errors: []error{},
	}
}

func (s *SocialMedia) UpdateBuilder() *SocialMediaBuilder {
	if s.id == 0 {
		return &SocialMediaBuilder{
			errors: []error{ErrInvalidID},
		}
	}
	if s.IsDeleted() {
		return &SocialMediaBuilder{
			errors: []error{ErrSocialMediaDeleted},
		}
	}
	return &SocialMediaBuilder{
		socialMedia: s,
		errors:      []error{},
	}
}

func (s *SocialMediaBuilder) WithID(id uint) *SocialMediaBuilder {
	s.socialMedia.id = id
	if id == 0 {
		s.errors = append(s.errors, ErrInvalidID)
	}
	return s
}

func (s *SocialMediaBuilder) WithIDCompany(id uint) *SocialMediaBuilder {
	s.socialMedia.idCompany = id
	if id == 0 {
		s.errors = append(s.errors, ErrInvalidIDCompany)
	}
	return s
}

func (s *SocialMediaBuilder) WithName(name string) *SocialMediaBuilder {
	s.socialMedia.name = vo.Text(name)
	if !s.socialMedia.name.IsValid() {
		s.errors = append(s.errors, ErrInvalidName)
	}
	return s
}

func (s *SocialMediaBuilder) WithURL(url string) *SocialMediaBuilder {
	s.socialMedia.url = vo.URL(url)
	if !s.socialMedia.url.IsValid() {
		s.errors = append(s.errors, ErrInvalidURL)
	}
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

func (s *SocialMediaBuilder) Build() (*SocialMedia, error) {
	if len(s.errors) > 0 {
		return nil, NewValidationErrors(s.errors)
	}
	if err := s.socialMedia.validate(); err != nil {
		return nil, err
	}
	return s.socialMedia, nil
}

func (s *SocialMediaBuilder) Apply() error {
	if len(s.errors) > 0 {
		return NewValidationErrors(s.errors)
	}
	if err := s.socialMedia.validate(); err != nil {
		return err
	}
	s.socialMedia.updatedAt = time.Now()
	return nil
}

func (s *SocialMedia) ID() uint              { return s.id }
func (s *SocialMedia) IDCompany() uint       { return s.idCompany }
func (s *SocialMedia) Name() vo.Text         { return s.name }
func (s *SocialMedia) URL() vo.URL           { return s.url }
func (s *SocialMedia) CreatedAt() time.Time  { return s.createdAt }
func (s *SocialMedia) UpdatedAt() time.Time  { return s.updatedAt }
func (s *SocialMedia) DeletedAt() *time.Time { return s.deletedAt }

func (s *SocialMedia) Delete() error {
	if s.id == 0 {
		return ErrInvalidID
	}
	if s.IsDeleted() {
		return ErrSocialMediaDeleted
	}
	now := time.Now()
	s.deletedAt = &now
	s.updatedAt = now
	return nil
}

func (s *SocialMedia) IsDeleted() bool {
	return s.deletedAt != nil
}

func (s *SocialMedia) SetID(id uint) error {
	if s.id != 0 {
		return ErrInvalidID
	}
	s.id = id
	return nil
}

func (s *SocialMedia) SetCompanyID(id uint) error {
	if s.idCompany != 0 {
		return ErrInvalidIDCompany
	}
	s.idCompany = id
	return nil
}

func (s *SocialMedia) validate() error {
	if !s.name.IsValid() {
		return ErrInvalidName
	}
	if !s.url.IsValid() {
		return ErrInvalidURL
	}
	return nil
}
