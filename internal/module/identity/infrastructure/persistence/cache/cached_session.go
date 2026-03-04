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
	repo               postgres.Session
	client             redis.UniversalClient
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
	repo postgres.Session,
	client redis.UniversalClient,
	logger repository.Logger,
	opts ...Option,
) CachedSession {
	c := CachedSession{
		repo,
		client,
		logger,
		"",
		true,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c CachedSession) kSession(sid int64) string {
	return c.prefix + "session:" + strconv.FormatInt(sid, 10)
}

func (c CachedSession) kUserSessions(uid int64) string {
	return c.prefix + "user_sessions:" + strconv.FormatInt(uid, 10)
}

func encodeSessionCache(uid int64, exp time.Time) string {
	return strconv.FormatInt(uid, 10) + "|" + strconv.FormatInt(exp.Unix(), 10)
}

func decodeSessionCache(v string) (uid int64, exp time.Time, ok bool) {
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
	return u, time.Unix(e, 0).UTC(), true
}

func ttlUntil(now, exp time.Time) time.Duration {
	d := exp.Sub(now)
	if d < 0 {
		return 0
	}
	return d
}

func (c CachedSession) Create(ctx context.Context, s *domain.Session) error {
	if err := c.repo.Create(ctx, s); err != nil {
		return err
	}
	now := time.Now()
	ttl := ttlUntil(now, s.ExpiresAt())
	if ttl > 0 {
		key := c.kSession(s.ID().Value())
		val := encodeSessionCache(s.IDUser().Value(), s.ExpiresAt())
		pipe := c.client.Pipeline()
		pipe.Set(ctx, key, val, ttl)
		pipe.SAdd(ctx, c.kUserSessions(s.ID().Value()), s.ID().Value())
		_, _ = pipe.Exec(ctx)
	}
	return nil
}

func (c CachedSession) GetByID(ctx context.Context, sid int64) (*domain.Session, error) {
	key := c.kSession(sid)
	v, err := c.client.Get(ctx, key).Result()
	if err == nil {
		uid, exp, ok := decodeSessionCache(v)
		if ok {
			now := time.Now()
			if !now.Before(exp) {
				_ = c.client.Del(ctx, key).Err()
				c.logger.Info(ctx, "session cache expired",
					"sid", sid,
					"user_id", uid,
				)
			} else {
				c.logger.Info(ctx, "session cache hit",
					"sid", sid,
					"user_id", uid,
				)
				sess, err := c.repo.GetByID(ctx, sid)
				if err == nil {
					return sess, nil
				}
				_ = c.client.Del(ctx, key).Err()
				return nil, err
			}
		}
	}
	if err != nil && err != redis.Nil {
		c.logger.Error(ctx, "session cache error, falling back to postgres",
			"sid", sid,
			"error", err,
		)
	}
	if !c.fallbackToPostgres {
		return nil, domain.ErrSessionNotFound
	}
	sess, err := c.repo.GetByID(ctx, sid)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if sess.RevokedAt() == nil && now.Before(sess.ExpiresAt()) {
		ttl := ttlUntil(now, sess.ExpiresAt())
		if ttl > 0 {
			val := encodeSessionCache(sess.IDUser().Value(), sess.ExpiresAt())
			pipe := c.client.Pipeline()
			pipe.Set(ctx, key, val, ttl)
			pipe.SAdd(ctx, c.kUserSessions(sess.IDUser().Value()), sid)
			_, _ = pipe.Exec(ctx)
		}
	}
	return sess, nil
}

func (c CachedSession) RotateRefreshToken(ctx context.Context, sid int64, hash [32]byte, expiresAt time.Time) error {
	if err := c.repo.RotateRefreshToken(ctx, sid, hash, expiresAt); err != nil {
		return err
	}
	key := c.kSession(sid)
	v, err := c.client.Get(ctx, key).Result()
	if err == nil {
		uid, _, ok := decodeSessionCache(v)
		if ok {
			now := time.Now()
			ttl := ttlUntil(now, expiresAt)
			if ttl > 0 {
				pipe := c.client.Pipeline()
				pipe.Set(ctx, key, encodeSessionCache(uid, expiresAt), ttl)
				pipe.SAdd(ctx, c.kUserSessions(uid), sid)
				_, _ = pipe.Exec(ctx)
			} else {
				_, _ = c.client.Del(ctx, key).Result()
			}
		}
	}
	return nil
}

func (c CachedSession) Revoke(ctx context.Context, sid int64) error {
	if err := c.repo.Revoke(ctx, sid); err != nil {
		return err
	}
	_, _ = c.client.Del(ctx, c.kSession(sid)).Result()
	return nil
}

func (c CachedSession) RevokeAllByUser(ctx context.Context, uid int64) error {
	if err := c.repo.RevokeAllByUser(ctx, uid); err != nil {
		return err
	}
	setKey := c.kUserSessions(uid)
	sids, err := c.client.SMembers(ctx, setKey).Result()
	if err == nil && len(sids) > 0 {
		pipe := c.client.Pipeline()
		for _, sidStr := range sids {
			pipe.Del(ctx, c.prefix+"session:"+sidStr)
		}
		pipe.Del(ctx, setKey)
		_, _ = pipe.Exec(ctx)
	} else {
		_, _ = c.client.Del(ctx, setKey).Result()
	}
	return nil
}

func (c CachedSession) CountActiveByUser(ctx context.Context, uid int64) (int, error) {
	return c.repo.CountActiveByUser(ctx, uid)
}

func (c CachedSession) IsActive(ctx context.Context, sid int64, uid int64, now time.Time) (bool, error) {
	start := time.Now()
	key := c.kSession(sid)
	v, err := c.client.Get(ctx, key).Result()
	if err == nil {
		idUser, exp, ok := decodeSessionCache(v)
		if !ok || idUser != uid || !now.Before(exp) {
			_, _ = c.client.Del(ctx, key).Result()
			c.logger.Info(ctx, "session.is_active cache_stale",
				"sid", sid,
				"user_id", uid,
				"took", time.Since(start).String(),
			)
			return false, nil
		}
		c.logger.Info(ctx, "session.is_active cache_hit",
			"sid", sid,
			"user_id", uid,
			"took", time.Since(start).String(),
		)
		return true, nil
	}
	if err != redis.Nil {
		c.logger.Error(ctx, "session.is_active redis_error_fallback_pg",
			"sid", sid,
			"user_id", uid,
			"error", err,
			"took", time.Since(start).String(),
		)
	} else {
		c.logger.Info(ctx, "session.is_active cache_miss_fallback_pg",
			"sid", sid,
			"user_id", uid,
			"took", time.Since(start).String(),
		)
	}
	ok, err := c.repo.IsActive(ctx, sid, uid, now)
	if err != nil || !ok {
		return ok, err
	}
	sess, err := c.repo.GetByID(ctx, sid)
	if err == nil && sess.RevokedAt() == nil && now.Before(sess.ExpiresAt()) {
		ttl := ttlUntil(now, sess.ExpiresAt())
		if ttl > 0 {
			val := encodeSessionCache(sess.IDUser().Value(), sess.ExpiresAt())
			pipe := c.client.Pipeline()
			pipe.Set(ctx, key, val, ttl)
			pipe.SAdd(ctx, c.kUserSessions(sess.IDUser().Value()), sid)
			_, _ = pipe.Exec(ctx)
		}
	}
	return true, nil
}

func (c CachedSession) RotateRefreshTokenAtomic(
	ctx context.Context,
	sid int64,
	oldHash [32]byte,
	hash [32]byte,
	expiryAtTime time.Time,
	now time.Time,
) (bool, error) {
	rotated, err := c.repo.RotateRefreshTokenAtomic(
		ctx,
		sid,
		oldHash,
		hash,
		expiryAtTime,
		now,
	)
	if err != nil {
		return false, err
	}
	if !rotated {
		return false, nil
	}
	key := c.kSession(sid)
	v, err := c.client.Get(ctx, key).Result()
	if err == nil {
		uid, _, ok := decodeSessionCache(v)
		if ok {
			ttl := ttlUntil(now, expiryAtTime)
			if ttl > 0 {
				pipe := c.client.Pipeline()
				pipe.Set(ctx, key, encodeSessionCache(uid, expiryAtTime), ttl)
				pipe.SAdd(ctx, c.kUserSessions(uid), sid)
				_, _ = pipe.Exec(ctx)
			} else {
				_, _ = c.client.Del(ctx, key).Result()
			}
			return true, nil
		}
	}
	sess, err := c.repo.GetByID(ctx, sid)
	if err == nil && sess != nil && sess.RevokedAt() == nil && now.Before(sess.ExpiresAt()) {
		ttl := ttlUntil(now, sess.ExpiresAt())
		if ttl > 0 {
			pipe := c.client.Pipeline()
			pipe.Set(ctx, key, encodeSessionCache(sess.IDUser().Value(), sess.ExpiresAt()), ttl)
			pipe.SAdd(ctx, c.kUserSessions(sess.IDUser().Value()), sid)
			_, _ = pipe.Exec(ctx)
		}
	}
	return true, nil
}
