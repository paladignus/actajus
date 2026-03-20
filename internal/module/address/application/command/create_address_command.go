// Package command
package command

type CreateAddressCommand struct {
	ZIP          string
	Title        string
	Street       string
	Number       uint
	Complement   string
	Reference    string
	Neighborhood string
	City         string
	State        string
	Country      string
}
