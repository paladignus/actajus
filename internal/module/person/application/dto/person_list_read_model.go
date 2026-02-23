// Package dto
package dto

type PersonListReadModel struct {
	People []PersonReadModel `json:"people"`
}
