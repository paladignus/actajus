// Package cache
package cache

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/domain"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/shared/application/repository"
	"github.com/redis/go-redis/v9"
)

type CachedSession struct {
	pg                 postgres.Session
	rdb                redis.UniversalClient
	logger             repository.Logger
	prefix             string
	fallbackToPostgres bool
}

type Option func(CachedSession)

func WithPrefix(prefix string) Option {
	return func(c CachedSession) { c.prefix = prefix }
}

func WithFallbackToPostgres(v bool) Option {
	return func(c CachedSession) { c.fallbackToPostgres = v }
}

func NewCachedSession(
	pg postgres.Session,
	rdb redis.UniversalClient,
	logger repository.Logger,
	opts ...Option,
) CachedSession {
	c := CachedSession{
		pg,
		rdb,
		logger,
		"",
		true,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c CachedSession) kSession(sid domain.IDSession) string {
	return c.prefix + "session:" + strconv.FormatInt(sid.Value(), 10)
}

func (c CachedSession) kUserSessions(uid domain.IDUser) string {
	return c.prefix + "user_sessions:" + strconv.FormatInt(uid.Value(), 10)
}

func encodeSessionCache(uid domain.IDUser, exp time.Time) string {
	return strconv.FormatInt(uid.Value(), 10) + "|" + strconv.FormatInt(exp.Unix(), 10)
}

func decodeSessionCache(v string) (uid domain.IDUser, exp time.Time, ok bool) {
	parts := strings.Split(v, "|")
	if len(parts) != 2 {
		return 0, time.Time{}, false
	}
	u, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, time.Time{}, false
	}
	e, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, time.Time{}, false
	}
	return domain.IDUser(u), time.Unix(e, 0).UTC(), true
}

func ttlUntil(now, exp time.Time) time.Duration {
	d := exp.Sub(now)
	if d < 0 {
		return 0
	}
	return d
}

func (c CachedSession) Create(ctx context.Context, s *domain.Session) error {
	if err := c.pg.Create(ctx, s); err != nil {
		return err
	}
	now := time.Now()
	ttl := ttlUntil(now, s.ExpiresAt())
	if ttl > 0 {
		key := c.kSession(s.ID())
		val := encodeSessionCache(s.IDUser(), s.ExpiresAt())
		pipe := c.rdb.Pipeline()
		pipe.Set(ctx, key, val, ttl)
		pipe.SAdd(ctx, c.kUserSessions(s.IDUser()), s.ID().Value())
		_, _ = pipe.Exec(ctx)
	}
	return nil
}

func (c CachedSession) GetByID(ctx context.Context, sid domain.IDSession) (*domain.Session, error) {
	key := c.kSession(sid)

	// 1) Tentar cache Redis
	v, err := c.rdb.Get(ctx, key).Result()
	if err == nil {
		uid, exp, ok := decodeSessionCache(v)
		if ok {
			now := time.Now()
			if !now.Before(exp) {
				// Cache expirado, remover
				_ = c.rdb.Del(ctx, key).Err()
				c.logger.Info(ctx, "session cache expired",
					"sid", sid.Value(),
					"user_id", uid.Value(),
				)
			} else {
				// Cache hit válido
				c.logger.Info(ctx, "session cache hit",
					"sid", sid.Value(),
					"user_id", uid.Value(),
				)
				// Retornar sessão completa do Postgres para garantir dados completos
				// O cache armazena apenas dados mínimos (uid, exp)
				sess, err := c.pg.GetByID(ctx, sid)
				if err == nil {
					return sess, nil
				}
				// Se falhou no Postgres, remover do Redis
				_ = c.rdb.Del(ctx, key).Err()
				return nil, err
			}
		}
	}

	// 2) Cache miss ou erro no Redis
	if err != nil && err != redis.Nil {
		c.logger.Error(ctx, "session cache error, falling back to postgres",
			"sid", sid.Value(),
			"error", err,
		)
	}

	// 3) Fallback para Postgres
	if !c.fallbackToPostgres {
		return nil, domain.ErrSessionNotFound
	}

	sess, err := c.pg.GetByID(ctx, sid)
	if err != nil {
		return nil, err
	}

	// 4) Cacheia se sessão estiver ativa (não revogada e não expirada)
	now := time.Now()
	if sess.RevokedAt() == nil && now.Before(sess.ExpiresAt()) {
		ttl := ttlUntil(now, sess.ExpiresAt())
		if ttl > 0 {
			val := encodeSessionCache(sess.IDUser(), sess.ExpiresAt())
			pipe := c.rdb.Pipeline()
			pipe.Set(ctx, key, val, ttl)
			pipe.SAdd(ctx, c.kUserSessions(sess.IDUser()), sid.Value())
			_, _ = pipe.Exec(ctx)
		}
	}

	return sess, nil
}

func (c CachedSession) RotateRefreshToken(ctx context.Context, sid domain.IDSession, newHash [32]byte, newExpiresAt time.Time) error {
	if err := c.pg.RotateRefreshToken(ctx, sid, newHash, newExpiresAt); err != nil {
		return err
	}
	key := c.kSession(sid)
	v, err := c.rdb.Get(ctx, key).Result()
	if err == nil {
		uid, _, ok := decodeSessionCache(v)
		if ok {
			now := time.Now()
			ttl := ttlUntil(now, newExpiresAt)
			if ttl > 0 {
				pipe := c.rdb.Pipeline()
				pipe.Set(ctx, key, encodeSessionCache(uid, newExpiresAt), ttl)
				pipe.SAdd(ctx, c.kUserSessions(uid), sid.Value())
				_, _ = pipe.Exec(ctx)
			} else {
				_, _ = c.rdb.Del(ctx, key).Result()
			}
		}
	}
	return nil
}

func (c CachedSession) Revoke(ctx context.Context, sid domain.IDSession) error {
	if err := c.pg.Revoke(ctx, sid); err != nil {
		return err
	}
	_, _ = c.rdb.Del(ctx, c.kSession(sid)).Result()
	return nil
}

func (c CachedSession) RevokeAllByUser(ctx context.Context, userID domain.IDUser) error {
	if err := c.pg.RevokeAllByUser(ctx, userID); err != nil {
		return err
	}
	setKey := c.kUserSessions(userID)
	sids, err := c.rdb.SMembers(ctx, setKey).Result()
	if err == nil && len(sids) > 0 {
		pipe := c.rdb.Pipeline()
		for _, sidStr := range sids {
			pipe.Del(ctx, c.prefix+"session:"+sidStr)
		}
		pipe.Del(ctx, setKey)
		_, _ = pipe.Exec(ctx)
	} else {
		_, _ = c.rdb.Del(ctx, setKey).Result()
	}
	return nil
}

func (c CachedSession) CountActiveByUser(ctx context.Context, idUser domain.IDUser) (int, error) {
	// Para manter simples e correto: usa Postgres (source of truth).
	// Otimização futura: count via Redis set (mas precisa limpar entradas expiradas).
	return c.pg.CountActiveByUser(ctx, idUser)
}

func (c CachedSession) IsActive(ctx context.Context, sid domain.IDSession, idUser domain.IDUser, now time.Time) (bool, error) {
	start := time.Now()
	key := c.kSession(sid)
	v, err := c.rdb.Get(ctx, key).Result()
	if err == nil {
		uid, exp, ok := decodeSessionCache(v)
		if !ok || uid.Value() != idUser.Value() || !now.Before(exp) {
			_, _ = c.rdb.Del(ctx, key).Result()
			c.logger.Info(ctx, "session.is_active cache_stale",
				"sid", sid.Value(),
				"user_id", idUser.Value(),
				"took", time.Since(start).String(),
			)
			return false, nil
		}
		c.logger.Info(ctx, "session.is_active cache_hit",
			"sid", sid.Value(),
			"user_id", idUser.Value(),
			"took", time.Since(start).String(),
		)
		return true, nil
	}
	if err != redis.Nil {
		c.logger.Error(ctx, "session.is_active redis_error_fallback_pg",
			"sid", sid.Value(),
			"user_id", idUser.Value(),
			"error", err,
			"took", time.Since(start).String(),
		)
	} else {
		c.logger.Info(ctx, "session.is_active cache_miss_fallback_pg",
			"sid", sid.Value(),
			"user_id", idUser.Value(),
			"took", time.Since(start).String(),
		)
	}
	ok, err := c.pg.IsActive(ctx, sid, idUser, now)
	if err != nil || !ok {
		return ok, err
	}
	sess, err := c.pg.GetByID(ctx, sid)
	if err == nil && sess.RevokedAt() == nil && now.Before(sess.ExpiresAt()) {
		ttl := ttlUntil(now, sess.ExpiresAt())
		if ttl > 0 {
			val := encodeSessionCache(sess.IDUser(), sess.ExpiresAt())
			pipe := c.rdb.Pipeline()
			pipe.Set(ctx, key, val, ttl)
			pipe.SAdd(ctx, c.kUserSessions(sess.IDUser()), sid.Value())
			_, _ = pipe.Exec(ctx)
		}
	}
	return true, nil
}

func (c CachedSession) RotateRefreshTokenAtomic(
	ctx context.Context,
	sid domain.IDSession,
	expectedOldHash [32]byte,
	newHash [32]byte,
	newExpiryAtTime time.Time,
	now time.Time,
) (bool, error) {
	rotated, err := c.pg.RotateRefreshTokenAtomic(
		ctx,
		sid,
		expectedOldHash,
		newHash,
		newExpiryAtTime,
		now,
	)
	if err != nil {
		return false, err
	}
	if !rotated {
		return false, nil
	}
	key := c.kSession(sid)
	v, err := c.rdb.Get(ctx, key).Result()
	if err == nil {
		uid, _, ok := decodeSessionCache(v)
		if ok {
			ttl := ttlUntil(now, newExpiryAtTime)
			if ttl > 0 {
				pipe := c.rdb.Pipeline()
				pipe.Set(ctx, key, encodeSessionCache(uid, newExpiryAtTime), ttl)
				pipe.SAdd(ctx, c.kUserSessions(uid), sid.Value())
				_, _ = pipe.Exec(ctx)
			} else {
				_, _ = c.rdb.Del(ctx, key).Result()
			}
			return true, nil
		}
	}
	sess, err := c.pg.GetByID(ctx, sid)
	if err == nil && sess != nil && sess.RevokedAt() == nil && now.Before(sess.ExpiresAt()) {
		ttl := ttlUntil(now, sess.ExpiresAt())
		if ttl > 0 {
			pipe := c.rdb.Pipeline()
			pipe.Set(ctx, key, encodeSessionCache(sess.IDUser(), sess.ExpiresAt()), ttl)
			pipe.SAdd(ctx, c.kUserSessions(sess.IDUser()), sid.Value())
			_, _ = pipe.Exec(ctx)
		}
	}
	return true, nil
}
