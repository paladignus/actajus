// Package service
package service

import "context"

type RBACRoleUsersIndex interface {
	AddUserToRole(ctx context.Context, idRole int16, idUser int64)
	RemoveUserFromRole(ctx context.Context, idRole int16, idUser int64)
	ListUsersByRole(ctx context.Context, idRole int16) ([]int64, error)
	AddUsersToRole(ctx context.Context, idRole int16, idUsers []int64)
}
