// Package security
package security

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/domain"
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

func (s *AuthorizationService) kUserPerms(uid domain.IDUser) string {
	return s.prefix + "user_perms:" + itoa64(uid.Value())
}

// HasPermission: Redis SISMEMBER -> fallback Postgres ListPermissionsByUser -> repopula cache
func (s *AuthorizationService) HasPermission(ctx context.Context, userID domain.IDUser, perm string) (bool, error) {
	key := s.kUserPerms(userID)
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
	perms, err := s.repo.ListPermissionsByUser(ctx, userID)
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

// func (s *AuthorizationService) HasPermission(ctx context.Context, userID domain.IDUser, perm string) (bool, error) {
// 	key := s.kUserPerms(userID)
// 	ok, err := s.rdb.SIsMember(ctx, key, perm).Result()
// 	fmt.Println("HasPermission:", key, perm, ok, err)
// 	if err == nil {
// 		return ok, nil
// 	}
// 	// redis.Nil não é comum em SIsMember, mas redis pode falhar.
// 	// Se falhar, segue fallback no Postgres.
//
// 	perms, err2 := s.repo.ListPermissionsByUser(ctx, userID)
// 	if err2 != nil {
// 		return false, err2
// 	}
// 	// repopula cache (best-effort)
// 	if len(perms) > 0 {
// 		members := make([]any, 0, len(perms))
// 		for _, p := range perms {
// 			members = append(members, p)
// 		}
// 		pipe := s.rdb.Pipeline()
// 		pipe.Del(ctx, key)
// 		pipe.SAdd(ctx, key, members...)
// 		pipe.Expire(ctx, key, s.ttl)
// 		_, _ = pipe.Exec(ctx)
// 	} else {
// 		_ = s.rdb.Expire(ctx, key, s.ttl).Err()
// 	}
//
// 	for _, p := range perms {
// 		if p == perm {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

func (s *AuthorizationService) InvalidateUser(ctx context.Context, userID domain.IDUser) {
	_, _ = s.rdb.Del(ctx, s.kUserPerms(userID)).Result()
}

// helper sem fmt (hot path)
func itoa64(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var b [32]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
