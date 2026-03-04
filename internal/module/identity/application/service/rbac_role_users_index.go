// Package service
package service

import "context"

type RBACRoleUsersIndex interface {
	AddUserToRole(ctx context.Context, rid int16, uid int64)
	RemoveUserFromRole(ctx context.Context, rid int16, uid int64)
	ListUsersByRole(ctx context.Context, rid int16) ([]int64, error)
	AddUsersToRole(ctx context.Context, rid int16, uid []int64)
}
