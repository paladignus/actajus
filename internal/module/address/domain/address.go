// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
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
}

func NewAddressBuilder() *AddressBuilder {
	now := time.Now()
	return &AddressBuilder{
		address: &Address{
			createdAt: now,
			updatedAt: now,
		},
	}
}

// func (a *Address) UpdateBuilder() *AddressBuilder {
// 	return &AddressBuilder{
// 		address: a,
// 	}
// }

func (b *AddressBuilder) WithID(id uint) *AddressBuilder {
	b.address.id = id
	return b
}

func (b *AddressBuilder) WithZIP(zip string) *AddressBuilder {
	b.address.zip = vo.ZIP(zip)
	return b
}

func (b *AddressBuilder) WithTitle(title string) *AddressBuilder {
	b.address.title = vo.Text(title)
	return b
}

func (b *AddressBuilder) WithStreet(street string) *AddressBuilder {
	b.address.street = vo.Text(street)
	return b
}

func (b *AddressBuilder) WithNumber(number uint) *AddressBuilder {
	b.address.number = number
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
	return b
}

func (b *AddressBuilder) WithCity(city string) *AddressBuilder {
	b.address.city = vo.Text(city)
	return b
}

func (b *AddressBuilder) WithState(state string) *AddressBuilder {
	b.address.state = vo.Text(state)
	return b
}

func (b *AddressBuilder) WithCountry(country string) *AddressBuilder {
	b.address.country = vo.Text(country)
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
	if err := b.address.validate(); err != nil {
		return nil, err
	}
	return b.address, nil
}

// func (b *AddressBuilder) Apply() error {
// 	if err := b.address.validate(); err != nil {
// 		return err
// 	}
// 	b.address.updatedAt = time.Now()
// 	return nil
// }

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
		return domain.NewFieldError("id", "address ID is invalid")
	}
	if a.IsDeleted() {
		return domain.NewFieldError("deleted_at", "address is already deleted")
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
		return domain.NewFieldError("id", "address ID is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "address ID is invalid")
	}
	a.id = id
	return nil
}

func (a *Address) validate() error {
	if a.zip != "" && !a.zip.IsValid() {
		return domain.NewFieldError("zip", "address ZIP is invalid")
	}
	if a.title != "" && !a.title.IsValid() {
		return domain.NewFieldError("title", "address title is invalid")
	}
	if a.street != "" && !a.street.IsValid() {
		return domain.NewFieldError("street", "address street is invalid")
	}
	if a.number == 0 {
		return domain.NewFieldError("number", "address number is invalid")
	}
	if a.neighborhood != "" && !a.neighborhood.IsValid() {
		return domain.NewFieldError("neighborhood", "address neighborhood is invalid")
	}
	if a.city != "" && !a.city.IsValid() {
		return domain.NewFieldError("city", "address city is invalid")
	}
	if a.state != "" && !a.state.IsValid() {
		return domain.NewFieldError("state", "address state is invalid")
	}
	if a.country != "" && !a.country.IsValid() {
		return domain.NewFieldError("country", "address country is invalid")
	}
	return nil
}
