package authx

import (
	"context"
	"github.com/shrinex/shield/codec"
	"github.com/shrinex/shield/semgt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strings"
	"time"
)

type (
	RedisSessionRepository struct {
		*redis.Redis
		flushMode   FlushMode
		codec       codec.Codec
		timeout     time.Duration
		idleTimeout time.Duration
	}
)

var _ semgt.Repository[*RedisSession] = (*RedisSessionRepository)(nil)

func NewRepository(conf redis.RedisConf, flushMode FlushMode, codec codec.Codec,
	timeout time.Duration, idleTimeout time.Duration) *RedisSessionRepository {
	return &RedisSessionRepository{
		codec:       codec,
		timeout:     timeout,
		flushMode:   flushMode,
		idleTimeout: idleTimeout,
		Redis:       redis.MustNewRedis(conf),
	}
}

func (r *RedisSessionRepository) Save(ctx context.Context, session *RedisSession) error {
	return session.flushNow(ctx)
}

func (r *RedisSessionRepository) Remove(ctx context.Context, token string) error {
	session, err := r.readSession(ctx, token, true)
	if err != nil {
		return err
	}

	if session == nil {
		return nil
	}

	err = session.setIdleTimeout(ctx, time.Duration(0))
	if err != nil {
		return err
	}

	return session.flushNow(ctx)
}

func (r *RedisSessionRepository) Create(ctx context.Context, token string) (*RedisSession, error) {
	cached := semgt.NewSessionTimeout(token, r.codec, r.timeout, r.idleTimeout)
	result := newSession(cached, r.Redis, r.flushMode, true)
	err := result.flushIfRequired(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RedisSessionRepository) Read(ctx context.Context, token string) (*RedisSession, error) {
	return r.readSession(ctx, token, false)
}

func (r *RedisSessionRepository) readSession(ctx context.Context, token string, allowExpired bool) (*RedisSession, error) {
	entries, err := r.HgetallCtx(ctx, sessionKey(token))
	if err != nil || len(entries) == 0 {
		return nil, err
	}

	session := r.loadSession(token, entries)
	expired, err := session.Expired(ctx)
	if err != nil {
		return nil, err
	}

	if expired && !allowExpired {
		return nil, nil
	}

	return newSession(session, r.Redis, r.flushMode, false), nil
}

func (r *RedisSessionRepository) loadSession(token string, entries map[string]string) *semgt.MapSession {
	session := semgt.NewSession(token, r.codec)

	for key, value := range entries {
		if key == startTimeKey {
			startTimeMillis := parseInt(value)
			session.SetStartTime(time.UnixMilli(startTimeMillis))
		} else if key == timeoutKey {
			timeoutNanos := parseInt(value)
			session.SetTimeout(time.Duration(timeoutNanos))
		} else if key == idleTimeoutKey {
			idleTimeoutNanos := parseInt(value)
			session.SetIdleTimeout(time.Duration(idleTimeoutNanos))
		} else if key == lastAccessTimeKey {
			lastAccessTimeMillis := parseInt(value)
			session.SetLastAccessTime(time.UnixMilli(lastAccessTimeMillis))
		} else if strings.HasPrefix(key, sessionAttributeKeyPrefix) {
			session.SetRawAttribute(strings.TrimPrefix(key, sessionAttributeKeyPrefix+":"), value)
		}
	}

	return session
}
