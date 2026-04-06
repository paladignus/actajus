// Package entity
package entity

import (
	"time"

	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type Address struct {
	IDAddress    uint
	Number       uint
	Zip          vo.ZIP
	Title        vo.Text
	Street       vo.Text
	Complement   vo.Text
	Reference    vo.Text
	Neighborhood vo.Text
	City         vo.Text
	State        vo.Text
	Country      vo.Text
	DeletedAt    *time.Time
}

func NewAddress(address readmodel.AddressReadModel) Address {
	return Address{
		IDAddress:    address.IDAddress,
		Zip:          vo.ZIP(address.Zip),
		Title:        vo.Text(address.Title),
		Street:       vo.Text(address.Street),
		Number:       address.Number,
		Complement:   vo.Text(address.Complement),
		Reference:    vo.Text(address.Reference),
		Neighborhood: vo.Text(address.Neighborhood),
		City:         vo.Text(address.City),
		State:        vo.Text(address.State),
		Country:      vo.Text(address.Country),
	}
}

func (a Address) Create() error {
	return a.validate()
}

func (a Address) Update() error {
	if a.IDAddress == 0 {
		return exception.ErrInvalidIDAddress
	}
	return a.validate()
}

func (a Address) IsDeleted() bool {
	return a.DeletedAt != nil
}

func (a Address) validate() error {
	if a.Zip.Value() != "" && !a.Zip.IsValid() {
		return exception.ErrInvalidZip
	}
	if a.Title.Value() != "" && !a.Title.IsValid() {
		return exception.ErrInvalidTitle
	}
	if a.Street.Value() != "" && !a.Street.IsValid() {
		return exception.ErrInvalidStreet
	}
	if a.Neighborhood.Value() != "" && !a.Neighborhood.IsValid() {
		return exception.ErrInvalidNeighborhood
	}
	if a.City.Value() != "" && !a.City.IsValid() {
		return exception.ErrInvalidCity
	}
	if a.State.Value() != "" && !a.State.IsValid() {
		return exception.ErrInvalidState
	}
	if a.Country.Value() != "" && !a.Country.IsValid() {
		return exception.ErrInvalidCountry
	}
	return nil
}
