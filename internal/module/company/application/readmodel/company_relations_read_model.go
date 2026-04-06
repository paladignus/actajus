// Package readmodel
package readmodel

import "time"

type CompanyAddressReadModel struct {
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

type CompanyPhoneReadModel struct {
	ID         int64
	Number     string
	Kind       string
	Department string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CompanyEmailReadModel struct {
	ID        int64
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CompanySocialMediaReadModel struct {
	ID        int64
	IDCompany int64
	Platform  string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
