// Package command
package command

type UpdateAddressCommand struct {
	IDAddress    int64
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
