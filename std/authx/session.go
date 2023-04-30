package authx

import (
	"context"
	"fmt"
	"github.com/shrinex/shield/semgt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"sync"
	"time"
)

type (
	RedisSession struct {
		cached    *semgt.MapSession
		delta     map[string]string
		redis     *redis.Redis
		mutex     sync.Mutex
		flushMode FlushMode
	}
)

var _ semgt.Session = (*RedisSession)(nil)

func newSession(cached *semgt.MapSession, redis *redis.Redis,
	flushMode FlushMode, isNew bool) *RedisSession {
	session := &RedisSession{
		redis:     redis,
		cached:    cached,
		flushMode: flushMode,
		delta:     make(map[string]string),
	}

	if isNew {
		session.delta[timeoutKey] = formatInt(cached.GetTimeout().Nanoseconds())
		session.delta[startTimeKey] = formatInt(cached.GetStartTime().UnixMilli())
		session.delta[idleTimeoutKey] = formatInt(cached.GetIdleTimeout().Nanoseconds())
		session.delta[lastAccessTimeKey] = formatInt(cached.GetLastAccessTime().UnixMilli())
	}

	return session
}

func (s *RedisSession) Token() string {
	return s.cached.Token()
}

func (s *RedisSession) StartTime(ctx context.Context) (time.Time, error) {
	return s.cached.StartTime(ctx)
}

func (s *RedisSession) Timeout(ctx context.Context) (time.Duration, error) {
	return s.cached.Timeout(ctx)
}

func (s *RedisSession) IdleTimeout(ctx context.Context) (time.Duration, error) {
	return s.cached.IdleTimeout(ctx)
}

func (s *RedisSession) LastAccessTime(ctx context.Context) (time.Time, error) {
	return s.cached.LastAccessTime(ctx)
}

func (s *RedisSession) Attribute(ctx context.Context, key string, ptr any) (bool, error) {
	return s.cached.Attribute(ctx, key, ptr)
}

func (s *RedisSession) AttributeAsInt(ctx context.Context, key string) (int64, bool, error) {
	return s.cached.AttributeAsInt(ctx, key)
}

func (s *RedisSession) AttributeAsBool(ctx context.Context, key string) (bool, bool, error) {
	return s.cached.AttributeAsBool(ctx, key)
}

func (s *RedisSession) AttributeAsFloat(ctx context.Context, key string) (float64, bool, error) {
	return s.cached.AttributeAsFloat(ctx, key)
}

func (s *RedisSession) AttributeAsString(ctx context.Context, key string) (string, bool, error) {
	return s.cached.AttributeAsString(ctx, key)
}

func (s *RedisSession) AttributeKeys(ctx context.Context) ([]string, error) {
	return s.cached.AttributeKeys(ctx)
}

func (s *RedisSession) SetAttribute(ctx context.Context, key string, value any) error {
	err := s.cached.SetAttribute(ctx, key, value)
	if err != nil {
		return err
	}

	rawValue, found := s.cached.RawAttribute(key)

	s.mutex.Lock()
	if found {
		s.delta[attributeKey(key)] = rawValue
	} else {
		delete(s.delta, attributeKey(key))
	}
	s.mutex.Unlock()

	if key == semgt.AlreadyExpiredKey ||
		key == semgt.AlreadyOverflowKey ||
		key == semgt.AlreadyReplacedKey {
		return s.flushNow(ctx)
	}

	return s.flushIfRequired(ctx)
}

func (s *RedisSession) RemoveAttribute(ctx context.Context, key string) error {
	err := s.cached.RemoveAttribute(ctx, key)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	delete(s.delta, attributeKey(key))
	s.mutex.Unlock()

	err = s.flushIfRequired(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisSession) Expired(ctx context.Context) (bool, error) {
	return s.cached.Expired(ctx)
}

func (s *RedisSession) Touch(ctx context.Context) error {
	lastAccessTime := nowFunc()
	s.cached.SetLastAccessTime(lastAccessTime)

	s.mutex.Lock()
	s.delta[lastAccessTimeKey] = formatInt(lastAccessTime.UnixMilli())
	s.mutex.Unlock()

	return s.flushIfRequired(ctx)
}

func (s *RedisSession) Flush(ctx context.Context) error {
	err := s.cached.Flush(ctx)
	if err != nil {
		return err
	}

	return s.flushNow(ctx)
}

func (s *RedisSession) Stop(ctx context.Context) error {
	return s.cached.Stop(ctx)
}

//=====================================
//		      Internal
//=====================================

func (s *RedisSession) setIdleTimeout(ctx context.Context, idleTimeout time.Duration) error {
	s.cached.SetIdleTimeout(idleTimeout)

	s.mutex.Lock()
	s.delta[idleTimeoutKey] = formatInt(idleTimeout.Nanoseconds())
	s.mutex.Unlock()

	return s.flushIfRequired(ctx)
}

func (s *RedisSession) getExpired() bool {
	return s.cached.GetExpired()
}

func (s *RedisSession) flushIfRequired(ctx context.Context) error {
	if s.flushMode != FlushModeImmediate {
		return nil
	}

	return s.flushNow(ctx)
}

func (s *RedisSession) flushNow(ctx context.Context) error {
	s.mutex.Lock()
	if len(s.delta) == 0 {
		s.mutex.Unlock()
		return nil
	}

	attrs := make(map[string]string)
	for key, value := range s.delta {
		attrs[key] = value
	}
	s.delta = make(map[string]string)
	s.mutex.Unlock()

	key := sessionKey(s.cached.Token())
	err := s.redis.HmsetCtx(ctx, key, attrs)
	if err != nil {
		return err
	}

	expiry := s.evalExpiry()
	err = s.redis.ExpireatCtx(ctx, key, expiry)
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisSession) evalExpiry() int64 {
	lastAccessTime := s.cached.GetLastAccessTime()
	idleTimeout := s.cached.GetIdleTimeout()
	startTime := s.cached.GetStartTime()
	timeout := s.cached.GetTimeout()

	var expirySeconds int64
	deadline := startTime.Add(timeout)
	expireTime := lastAccessTime.Add(idleTimeout)
	if deadline.Before(expireTime) {
		expirySeconds = int64(deadline.Second())
	} else {
		expirySeconds = int64(expireTime.Second())
	}

	// 3min after
	return expirySeconds + int64(time.Second*180)
}

func sessionKey(token string) string {
	return fmt.Sprintf("%s:%s", sessionKeyPrefix, token)
}

func attributeKey(name string) string {
	return fmt.Sprintf("%s:%s", sessionAttributeKeyPrefix, name)
}
