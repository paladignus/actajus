// Package repository
package repository

import "context"

type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type UnitOfWorkEnterprise interface {
	UnitOfWork
	Enterprise() IEnterprise
	Address() IAddress
}

type UnitOfWorkDefault interface {
	Enterprise() UnitOfWorkEnterprise
}
