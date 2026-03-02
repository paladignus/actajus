// Package security
package security

import (
	"context"
)

type RoleUserPair struct {
	IDRole int16
	IDUser int64
}

func RebuildRoleUsersIndex(ctx context.Context, idx *RBACRoleUsersIndex, pairs []RoleUserPair) error {
	grouped := make(map[int16][]int64, 32)
	for _, p := range pairs {
		grouped[p.IDRole] = append(grouped[p.IDRole], p.IDUser)
	}
	for idRole, idUsers := range grouped {
		idx.AddUsersToRole(ctx, idRole, idUsers)
	}
	return nil
}
