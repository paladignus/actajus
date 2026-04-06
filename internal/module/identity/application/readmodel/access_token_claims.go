// Package readmodel
package readmodel

import (
	"time"
)

type AccessTokenClaims struct {
	IDUser    int64
	IDSession int64
	Issuer    string
	Audience  string
	ExpiresAt time.Time
	IssuedAt  time.Time
}
