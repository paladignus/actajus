// Package readmodel
package readmodel

import "time"

type AddressReadModel struct {
	ID           int64
	ZIP          string
	Title        string
	Street       string
	Number       uint
	Complement   *string
	Reference    *string
	Neighborhood string
	City         string
	State        string
	Country      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
