// Package security
package security

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RBACRoleUsersIndex struct {
	rdb    redis.UniversalClient
	prefix string
}

type RoleUsersIndexOption func(*RBACRoleUsersIndex)

func WithRoleUsersIndexPrefix(prefix string) RoleUsersIndexOption {
	return func(i *RBACRoleUsersIndex) { i.prefix = prefix }
}

func NewRBACRoleUsersIndex(rdb redis.UniversalClient, opts ...RoleUsersIndexOption) *RBACRoleUsersIndex {
	idx := &RBACRoleUsersIndex{
		rdb:    rdb,
		prefix: "rbac:",
	}
	for _, opt := range opts {
		opt(idx)
	}
	return idx
}

func (i *RBACRoleUsersIndex) key(idRole int16) string {
	return i.prefix + "role_users:" + strconv.FormatInt(int64(idRole), 10)
}

func (i *RBACRoleUsersIndex) AddUserToRole(ctx context.Context, idRole int16, IDUser int64) {
	_ = i.rdb.SAdd(ctx, i.key(idRole), IDUser).Err()
}

func (i *RBACRoleUsersIndex) RemoveUserFromRole(ctx context.Context, idRole int16, IDUser int64) {
	_ = i.rdb.SRem(ctx, i.key(idRole), IDUser).Err()
}

func (i *RBACRoleUsersIndex) ListUsersByRole(ctx context.Context, idRole int16) ([]int64, error) {
	values, err := i.rdb.SMembers(ctx, i.key(idRole)).Result()
	if err != nil {
		return nil, err
	}
	out := make([]int64, 0, len(values))
	for _, v := range values {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		out = append(out, id)
	}
	return out, nil
}

func (i *RBACRoleUsersIndex) AddUsersToRole(ctx context.Context, idRole int16, IDUsers []int64) {
	if len(IDUsers) == 0 {
		return
	}
	members := make([]any, 0, len(IDUsers))
	for _, uid := range IDUsers {
		members = append(members, uid)
	}
	_ = i.rdb.SAdd(ctx, i.key(idRole), members...).Err()
}
