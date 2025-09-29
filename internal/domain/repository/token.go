// Package repository
package repository

import "time"

type Token interface {
	SaveRevokedToken(string, time.Time) error
	IsTokenRevoked(string) (bool, error)
	// ClearExpiredTokens() error
}
