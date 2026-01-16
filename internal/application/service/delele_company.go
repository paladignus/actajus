// Package service
package service

import "context"

type IDeleteCompany interface {
	Execute(ctx context.Context, id uint) error
}
