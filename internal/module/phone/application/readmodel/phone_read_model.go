// Package readmodel
package readmodel

import "time"

type PhoneReadModel struct {
	ID         int64
	Number     string
	Kind       string
	Department string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
