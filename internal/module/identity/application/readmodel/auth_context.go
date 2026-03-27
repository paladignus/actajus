// Package readmodel
package readmodel

import "time"

type AuthContext struct {
	IDUser    int64
	IDSession int64
	IssuedAt  time.Time
	ExpiresAt time.Time
}
