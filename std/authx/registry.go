package authx

import (
	"context"
	"fmt"
	"github.com/shrinex/shield/security"
	"github.com/shrinex/shield/semgt"
)

type (
	RedisSessionRegistry struct {
		repo *RedisSessionRepository
	}
)

var _ semgt.Registry[*RedisSession] = (*RedisSessionRegistry)(nil)

func NewRegistry(repo *RedisSessionRepository) *RedisSessionRegistry {
	return &RedisSessionRegistry{repo: repo}
}

func (r *RedisSessionRegistry) Register(ctx context.Context, principal string, session *RedisSession) error {
	signature, err := r.buildSignature(ctx, session)
	if err != nil {
		return err
	}

	key := principalKey(principal)
	return r.addSignature(ctx, key, signature)
}

func (r *RedisSessionRegistry) Deregister(ctx context.Context, principal string, session *RedisSession) error {
	signature, err := r.buildSignature(ctx, session)
	if err != nil {
		return err
	}

	key := principalKey(principal)
	return r.removeSignatures(ctx, key, signature)
}

func (r *RedisSessionRegistry) ActiveSessions(ctx context.Context, principal string) ([]*RedisSession, error) {
	key := principalKey(principal)
	signatures, err := r.readSignatures(ctx, key)
	if err != nil {
		return nil, err
	}

	active := make([]*RedisSession, 0)
	if len(signatures) == 0 {
		return active, r.removeKey(ctx, key)
	}

	var expiry int64
	inactive := make([]string, 0)
	for _, signature := range signatures {
		var session *RedisSession
		session, err = r.readSession(ctx, signature, true)
		if err != nil {
			return nil, err
		}

		// deregister
		if session == nil {
			inactive = append(inactive, signature)
		} else {
			if expired := session.getExpired(); expired {
				inactive = append(inactive, signature)
			} else {
				active = append(active, session)
				evaluated := session.evalExpiry()
				if evaluated > expiry {
					expiry = evaluated
				}
			}
		}
	}

	// means no active
	if expiry == 0 {
		return active, r.removeKey(ctx, key)
	}

	// remove inactive
	if len(inactive) > 0 {
		_ = r.removeSignatures(ctx, key, inactive...)
	}

	err = r.expireKeyAt(ctx, key, expiry)
	if err != nil {
		return nil, err
	}

	return active, nil
}

func (r *RedisSessionRegistry) KeepAlive(ctx context.Context, principal string) error {
	key := principalKey(principal)
	signatures, err := r.readSignatures(ctx, key)
	if err != nil {
		return err
	}

	if len(signatures) == 0 {
		return r.removeKey(ctx, key)
	}

	var expiry int64
	for _, signature := range signatures {
		var session *RedisSession
		session, err = r.readSession(ctx, signature, false)
		if err != nil {
			return err
		}

		if session != nil {
			evaluated := session.evalExpiry()
			if evaluated > expiry {
				expiry = evaluated
			}
		}
	}

	if expiry == 0 {
		return r.removeKey(ctx, key)
	}

	return r.expireKeyAt(ctx, key, expiry)
}

func (r *RedisSessionRegistry) buildSignature(ctx context.Context, session *RedisSession) (string, error) {
	platform, found, err := session.AttributeAsString(ctx, security.PlatformKey)
	if err != nil {
		return "", err
	}

	if !found || len(platform) == 0 {
		platform = security.DefaultPlatform
	}

	data := make(map[string]string)
	data[platformKey] = platform
	data[tokenKey] = session.Token()
	signature, err := r.repo.codec.Encode(data)
	if err != nil {
		return "", err
	}

	return signature, nil
}

func (r *RedisSessionRegistry) addSignature(ctx context.Context, key, signature string) error {
	_, err := r.repo.Redis.SaddCtx(ctx, key, signature)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisSessionRegistry) readSignatures(ctx context.Context, key string) ([]string, error) {
	signatures, err := r.repo.Redis.SmembersCtx(ctx, key)
	if err != nil {
		return nil, err
	}

	return signatures, nil
}

func (r *RedisSessionRegistry) removeSignatures(ctx context.Context, key string, signatures ...string) error {
	_, err := r.repo.Redis.SremCtx(ctx, key, signatures)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisSessionRegistry) removeKey(ctx context.Context, key string) error {
	_, err := r.repo.Redis.DelCtx(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisSessionRegistry) expireKeyAt(ctx context.Context, key string, expiry int64) error {
	return r.repo.ExpireatCtx(ctx, key, expiry)
}

func (r *RedisSessionRegistry) readSession(ctx context.Context, signature string, allowExpired bool) (*RedisSession, error) {
	data := make(map[string]string)
	err := r.repo.codec.Decode(signature, &data)
	if err != nil {
		return nil, err
	}

	return r.repo.readSession(ctx, data[tokenKey], allowExpired)
}

func principalKey(principal string) string {
	return fmt.Sprintf("%s:%s", sessionPrincipalKeyPrefix, principal)
}
