// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
)

func resolveRoleUsersWithFallback(
	ctx context.Context,
	roleID int16,
	index service.RBACRoleUsersIndex,
	repo repository.RoleUserAdminRepository,
) ([]int64, error) {
	userIDs, err := index.ListUsersByRole(ctx, roleID)
	if err == nil && len(userIDs) > 0 {
		return userIDs, nil
	}

	dbUsers, err := repo.ListUserIDsByRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	out := make([]int64, 0, len(dbUsers))
	for _, uid := range dbUsers {
		out = append(out, uid.Value())
	}
	if len(out) > 0 {
		index.AddUsersToRole(ctx, roleID, out)
	}
	return out, nil
}
