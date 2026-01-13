// Package entity
package entity

import (
	"github.com/paladignus/actajus/internal/application/dto"
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
	Neighborhood vo.Text
	City         vo.Text
	State        vo.Text
	Country      vo.Text
}

func NewAddress(address dto.Address) (Address, error) {
	a := Address{
		IDAddress:    address.IDAddress,
		Zip:          vo.ZIP(address.Zip),
		Title:        vo.Text(address.Title),
		Street:       vo.Text(address.Street),
		Number:       address.Number,
		Complement:   vo.Text(address.Complement),
		Neighborhood: vo.Text(address.Neighborhood),
		City:         vo.Text(address.City),
		State:        vo.Text(address.State),
		Country:      vo.Text(address.Country),
	}
	if !a.Zip.IsValid() {
		return a, exception.ErrInvalidZip
	}
	if !a.Title.IsValid() {
		return a, exception.ErrInvalidTitle
	}
	if !a.Street.IsValid() {
		return a, exception.ErrInvalidStreet
	}
	if !a.Neighborhood.IsValid() {
		return a, exception.ErrInvalidNeighborhood
	}
	if !a.City.IsValid() {
		return a, exception.ErrInvalidCity
	}
	if !a.State.IsValid() {
		return a, exception.ErrInvalidState
	}
	if a.Country.Value() != "" && !a.Country.IsValid() {
		return a, exception.ErrInvalidCountry
	}
	return a, nil
}
