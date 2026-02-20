// Package domain
package domain

import "context"

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}
