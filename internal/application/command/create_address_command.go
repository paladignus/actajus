// Package command
package command

type CreateAddressCommand struct {
	Number       uint
	Zip          string
	Title        string
	Street       string
	Complement   string
	Reference    string
	Neighborhood string
	City         string
	State        string
	Country      string
}
