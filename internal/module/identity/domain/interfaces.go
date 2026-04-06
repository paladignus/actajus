// Package domain provides domain entities and interfaces
package domain

import "time"

// PasswordHasher hashes and verifies passwords
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

// Clock provides current time
type Clock interface {
	Now() time.Time
}
