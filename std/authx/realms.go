package authx

import (
	"context"
	"github.com/shrinex/shield/authc"
	"github.com/shrinex/shield/security"
)

type (
	bearerAuthRealm struct {
		repo *RedisSessionRepository
	}
)

var _ authc.Realm = (*bearerAuthRealm)(nil)

func NewBearerAuthRealm(repo *RedisSessionRepository) authc.Realm {
	return &bearerAuthRealm{repo: repo}
}

func (r *bearerAuthRealm) Supports(token authc.Token) bool {
	_, ok := token.(*authc.BearerToken)
	return ok
}

func (r *bearerAuthRealm) LoadUserDetails(ctx context.Context, token authc.Token) (authc.UserDetails, error) {
	session, err := r.repo.Read(ctx, token.Principal())
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, authc.ErrUnauthenticated
	}

	var userDetails UserDetails
	found, err := session.Attribute(ctx, security.UserDetailsKey, &userDetails)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, authc.ErrUnauthenticated
	}

	return &userDetails, nil
}
