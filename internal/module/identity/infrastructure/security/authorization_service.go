// Package security
package security

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/paladignus/actajus/internal/module/identity/application/repository"
)

type AuthorizationService struct {
	repo   repository.Authorization
	rdb    redis.UniversalClient
	prefix string
	ttl    time.Duration
}

type AuthzOption func(*AuthorizationService)

func WithAuthzPrefix(p string) AuthzOption { return func(a *AuthorizationService) { a.prefix = p } }
func WithAuthzTTL(ttl time.Duration) AuthzOption {
	return func(a *AuthorizationService) { a.ttl = ttl }
}

func NewAuthorizationService(repo repository.Authorization, rdb redis.UniversalClient, opts ...AuthzOption) *AuthorizationService {
	s := &AuthorizationService{
		repo:   repo,
		rdb:    rdb,
		prefix: "rbac:",
		ttl:    10 * time.Minute,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *AuthorizationService) kUserPerms(uid int64) string {
	return s.prefix + "user_perms:" + strconv.FormatInt(uid, 10)
	// return fmt.Sprintf("%suser_perms:%d", s.prefix, uid)
	// return s.prefix + "user_perms:" + itoa64(uid)
}

// HasPermission: Redis SISMEMBER -> fallback Postgres ListPermissionsByUser -> repopula cache
func (s *AuthorizationService) HasPermission(ctx context.Context, uid int64, perm string) (bool, error) {
	key := s.kUserPerms(uid)
	// 1) Se cache existe, checa membership
	exists, err := s.rdb.Exists(ctx, key).Result()
	if err == nil && exists == 1 {
		ok, err := s.rdb.SIsMember(ctx, key, perm).Result()
		if err == nil {
			return ok, nil
		}
		// se deu erro no redis aqui, cai pro fallback
	}
	// 2) Cache miss (key inexistente) OU redis falhou: fallback Postgres
	perms, err := s.repo.ListPermissionsByUser(ctx, uid)
	if err != nil {
		return false, err
	}
	// 3) Popular cache (best-effort)
	if len(perms) > 0 {
		members := make([]any, 0, len(perms))
		for _, p := range perms {
			members = append(members, p)
		}
		pipe := s.rdb.Pipeline()
		pipe.SAdd(ctx, key, members...)
		pipe.Expire(ctx, key, s.ttl)
		_, _ = pipe.Exec(ctx)
	} else {
		// opcional: criar key vazia com TTL (pra evitar bater no DB repetidamente)
		// um jeito simples: SETEX marcador
		_ = s.rdb.Set(ctx, key+":empty", "1", s.ttl).Err()
	}
	// 4) Decide em memória
	for _, p := range perms {
		if p == perm {
			return true, nil
		}
	}
	return false, nil
}

func (s *AuthorizationService) InvalidateUser(ctx context.Context, uid int64) {
	_, _ = s.rdb.Del(ctx, s.kUserPerms(uid)).Result()
}

// helper sem fmt (hot path)
// func itoa64(v int64) string {
// 	if v == 0 {
// 		return "0"
// 	}
// 	neg := v < 0
// 	if neg {
// 		v = -v
// 	}
// 	var b [32]byte
// 	i := len(b)
// 	for v > 0 {
// 		i--
// 		b[i] = byte('0' + v%10)
// 		v /= 10
// 	}
// 	if neg {
// 		i--
// 		b[i] = '-'
// 	}
// 	return string(b[i:])
// }
