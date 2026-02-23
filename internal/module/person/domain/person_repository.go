// Package domain
package domain

import "context"

type PersonRepository interface {
	Create(ctx context.Context, person *Person) error
}
