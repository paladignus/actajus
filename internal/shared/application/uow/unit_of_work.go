// Package uow
package uow

import "context"

type Tx any

type UnitOfWork interface {
	Do(ctx context.Context, fn func(tx Tx) error) error
}
