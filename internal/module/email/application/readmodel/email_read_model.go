// Package readmodel
package readmodel

import "time"

type EmailReadModel struct {
	ID        int64
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
