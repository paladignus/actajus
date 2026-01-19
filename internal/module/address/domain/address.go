// Package domain
package domain

import (
	"time"

	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type Address struct {
	id           uint
	number       uint
	zip          vo.ZIP
	title        vo.Text
	street       vo.Text
	complement   vo.Text
	reference    vo.Text
	neighborhood vo.Text
	city         vo.Text
	state        vo.Text
	country      vo.Text
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    *time.Time
}

type AddressBuilder struct {
	address *Address
	errors  []error
}

func NewAddressBuilder() *AddressBuilder {
	now := time.Now()
	return &AddressBuilder{
		address: &Address{
			createdAt: now,
			updatedAt: now,
		},
		errors: []error{},
	}
}

func (a *Address) UpdateBuilder() *AddressBuilder {
	if a.id == 0 {
		return &AddressBuilder{
			errors: []error{ErrInvalidID},
		}
	}
	if a.IsDeleted() {
		return &AddressBuilder{
			errors: []error{ErrAddressDeleted},
		}
	}
	return &AddressBuilder{
		address: a,
		errors:  []error{},
	}
}

func (b *AddressBuilder) WithID(id uint) *AddressBuilder {
	b.address.id = id
	if id == 0 {
		b.errors = append(b.errors, ErrInvalidID)
	}
	return b
}

func (b *AddressBuilder) WithZIP(zip string) *AddressBuilder {
	b.address.zip = vo.ZIP(zip)
	if !b.address.zip.IsValid() {
		b.errors = append(b.errors, ErrInvalidZIP)
	}
	return b
}

func (b *AddressBuilder) WithTitle(title string) *AddressBuilder {
	b.address.title = vo.Text(title)
	if !b.address.title.IsValid() {
		b.errors = append(b.errors, ErrInvalidTitle)
	}
	return b
}

func (b *AddressBuilder) WithStreet(street string) *AddressBuilder {
	b.address.street = vo.Text(street)
	if !b.address.street.IsValid() {
		b.errors = append(b.errors, ErrInvalidStreet)
	}
	return b
}

func (b *AddressBuilder) WithNumber(number uint) *AddressBuilder {
	b.address.number = number
	if number == 0 {
		b.errors = append(b.errors, ErrInvalidNumber)
	}
	return b
}

func (b *AddressBuilder) WithComplement(complement string) *AddressBuilder {
	b.address.complement = vo.Text(complement)
	return b
}

func (b *AddressBuilder) WithReference(reference string) *AddressBuilder {
	b.address.reference = vo.Text(reference)
	return b
}

func (b *AddressBuilder) WithNeighborhood(neighborhood string) *AddressBuilder {
	b.address.neighborhood = vo.Text(neighborhood)
	if !b.address.neighborhood.IsValid() {
		b.errors = append(b.errors, ErrInvalidNeighborhood)
	}
	return b
}

func (b *AddressBuilder) WithCity(city string) *AddressBuilder {
	b.address.city = vo.Text(city)
	if !b.address.city.IsValid() {
		b.errors = append(b.errors, ErrInvalidCity)
	}
	return b
}

func (b *AddressBuilder) WithState(state string) *AddressBuilder {
	b.address.state = vo.Text(state)
	if !b.address.state.IsValid() {
		b.errors = append(b.errors, ErrInvalidState)
	}
	return b
}

func (b *AddressBuilder) WithCountry(country string) *AddressBuilder {
	b.address.country = vo.Text(country)
	if !b.address.country.IsValid() {
		b.errors = append(b.errors, ErrInvalidCountry)
	}
	return b
}

func (b *AddressBuilder) WithCreatedAt(createdAt time.Time) *AddressBuilder {
	b.address.createdAt = createdAt
	return b
}

func (b *AddressBuilder) WithUpdatedAt(updatedAt time.Time) *AddressBuilder {
	b.address.updatedAt = updatedAt
	return b
}

func (b *AddressBuilder) WithDeletedAt(deletedAt *time.Time) *AddressBuilder {
	b.address.deletedAt = deletedAt
	return b
}

func (b *AddressBuilder) Build() (*Address, error) {
	if len(b.errors) > 0 {
		return nil, NewValidationErrors(b.errors)
	}
	if err := b.address.validate(); err != nil {
		return nil, err
	}
	return b.address, nil
}

func (b *AddressBuilder) Apply() error {
	if len(b.errors) > 0 {
		return NewValidationErrors(b.errors)
	}
	if err := b.address.validate(); err != nil {
		return err
	}
	b.address.updatedAt = time.Now()
	return nil
}

func (a *Address) ID() uint              { return a.id }
func (a *Address) Number() uint          { return a.number }
func (a *Address) ZIP() vo.ZIP           { return a.zip }
func (a *Address) Title() vo.Text        { return a.title }
func (a *Address) Street() vo.Text       { return a.street }
func (a *Address) Complement() vo.Text   { return a.complement }
func (a *Address) Reference() vo.Text    { return a.reference }
func (a *Address) Neighborhood() vo.Text { return a.neighborhood }
func (a *Address) City() vo.Text         { return a.city }
func (a *Address) State() vo.Text        { return a.state }
func (a *Address) Country() vo.Text      { return a.country }
func (a *Address) CreatedAt() time.Time  { return a.createdAt }
func (a *Address) UpdatedAt() time.Time  { return a.updatedAt }
func (a *Address) DeletedAt() *time.Time { return a.deletedAt }

func (a *Address) Delete() error {
	if a.id == 0 {
		return ErrInvalidID
	}
	if a.IsDeleted() {
		return ErrAlreadyDeleted
	}
	now := time.Now()
	a.deletedAt = &now
	a.updatedAt = now
	return nil
}

func (a *Address) IsDeleted() bool {
	return a.deletedAt != nil
}

func (a *Address) SetID(id uint) error {
	if a.id != 0 {
		return ErrIDAlreadySet
	}
	if id == 0 {
		return ErrInvalidID
	}
	a.id = id
	return nil
}

func (a *Address) validate() error {
	if !a.zip.IsValid() {
		return ErrInvalidZIP
	}
	if !a.title.IsValid() {
		return ErrInvalidTitle
	}
	if !a.street.IsValid() {
		return ErrInvalidStreet
	}
	if a.number == 0 {
		return ErrInvalidNumber
	}
	if !a.neighborhood.IsValid() {
		return ErrInvalidNeighborhood
	}
	if !a.city.IsValid() {
		return ErrInvalidCity
	}
	if !a.state.IsValid() {
		return ErrInvalidState
	}
	if !a.country.IsValid() {
		return ErrInvalidCountry
	}
	return nil
}
