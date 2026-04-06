// Package repository
package repository

import "context"

type RoleAdminRepository interface {
	Create(ctx context.Context, name, description string) error
	Update(ctx context.Context, id int16, name, description string) error
	Delete(ctx context.Context, id int16) error
}

type PermissionAdminRepository interface {
	Create(ctx context.Context, resource, action, description string) error
	Update(ctx context.Context, id int16, resource, action, description string) error
	Delete(ctx context.Context, id int16) error
}
