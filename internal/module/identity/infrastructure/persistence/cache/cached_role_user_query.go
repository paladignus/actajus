// Package cache
package cache

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	sharedrepo "github.com/paladignus/actajus/internal/shared/application/repository"
)

// CachedRoleUserQueryRepository implementa RoleUserQueryRepository com cache Redis
// e fallback para o repositório Postgres.
//
// Padrão: Cache-Aside (Lazy Loading) com fallback automático.
// - Primeiro tenta obter do cache (RBACRoleUsersIndex)
// - Em caso de miss ou erro, busca no banco (RoleUserAdminRepository)
// - Repopula o cache após leitura do banco
type CachedRoleUserQueryRepository struct {
	repo   repository.RoleUserAdminRepository
	index  service.RBACRoleUsersIndex
	logger sharedrepo.Logger
}

// NewCachedRoleUserQueryRepository cria uma instância com cache e fallback.
func NewCachedRoleUserQueryRepository(
	repo repository.RoleUserAdminRepository,
	index service.RBACRoleUsersIndex,
	logger sharedrepo.Logger,
) repository.RoleUserQueryRepository {
	return &CachedRoleUserQueryRepository{
		repo:   repo,
		index:  index,
		logger: logger,
	}
}

// ListUserIDsByRole retorna os usuários de uma role usando cache com fallback.
//
// Fluxo:
// 1. Tenta obter do cache (index Redis)
// 2. Se cache miss ou erro, busca no banco (repo Postgres)
// 3. Repopula o cache se houver dados
// 4. Retorna os dados
func (r *CachedRoleUserQueryRepository) ListUserIDsByRole(
	ctx context.Context,
	rid int16,
) ([]int64, error) {
	// 1) Tentar cache primeiro (leitura rápida)
	uid, err := r.index.ListUsersByRole(ctx, rid)
	if err == nil && len(uid) > 0 {
		r.logger.Debug(ctx, "RBAC role users cache hit",
			"role_id", rid,
			"count", len(uid),
		)
		return uid, nil
	}
	// 2) Cache miss ou erro - log para observabilidade
	r.logger.Debug(ctx, "RBAC role users cache miss, fallback to database",
		"role_id", rid,
		"cache_error", err,
	)
	// 3) Fallback: busca no banco de dados
	uid, err = r.repo.ListUserIDsByRole(ctx, rid)
	if err != nil {
		r.logger.Error(ctx, "RBAC role users database query failed",
			"role_id", rid,
			"error", err,
		)
		return nil, err
	}
	// 4) Repopula cache se houver dados
	if len(uid) > 0 {
		// idUsers := r.toInt64Slice(dbUsers)
		r.index.AddUsersToRole(ctx, rid, uid)
		r.logger.Debug(ctx, "RBAC role users cache repopulated",
			"role_id", rid,
			"count", len(uid),
		)
	}
	return uid, nil
}

// toIDUsers converte []int64 para []domain.IDUser
// func (r *CachedRoleUserQueryRepository) toIDUsers(ids []int64) []int64 {
// 	out := make([]int64, 0, len(ids))
// 	for _, id := range ids {
// 		out = append(out, id)
// 	}
// 	return out
// }

// toInt64Slice converte []domain.IDUser para []int64
// func (r *CachedRoleUserQueryRepository) toInt64Slice(ids []int64) []int64 {
// 	out := make([]int64, 0, len(ids))
// 	for _, id := range ids {
// 		out = append(out, id)
// 	}
// 	return out
// }
