// Package event
package event

type DomainEvent interface {
	Subject() string
	Serialize() ([]byte, error)
}
