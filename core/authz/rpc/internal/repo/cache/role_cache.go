package cache

import (
	"context"
	"core/authz/proto/model"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"shrine/std/globals"
	"shrine/std/utils/slices"
)

type (
	RoleCache struct {
		cache cache.Cache
		dao   model.RoleModel
	}
)

func NewRoleCache(conf cache.CacheConf, dao model.RoleModel, opts ...cache.Option) *RoleCache {
	cc := cache.New(conf, syncx.NewSingleFlight(), cache.NewStat("resource cache"), ErrNotFound, opts...)
	return &RoleCache{dao: dao, cache: cc}
}

func (c *RoleCache) ListRoles(ctx context.Context, sysType, isAdmin, userId int64) (roles []*model.Role, err error) {
	if isAdmin == globals.FlagTrue {
		return c.ListRolesBySysType(ctx, sysType)
	}

	resp := slices.Empty[*model.Role]()
	key := fmt.Sprintf("%s:%d:%d", globals.RolesCacheKeyPrefix, sysType, userId)
	err = c.cache.TakeCtx(ctx, &resp, key, func(val any) error {
		roles, err = c.dao.ListRolesBySysTypeAndUserId(ctx, sysType, userId)
		if err != nil {
			return err
		}

		*(val.(*[]*model.Role)) = roles
		return nil
	})

	return resp, err
}

func (c *RoleCache) ListRolesBySysType(ctx context.Context, sysType int64) (roles []*model.Role, err error) {
	resp := slices.Empty[*model.Role]()
	key := fmt.Sprintf("%s:%d", globals.RolesCacheKeyPrefix, sysType)
	err = c.cache.TakeCtx(ctx, &resp, key, func(val any) error {
		roles, err = c.dao.ListAllRolesBySysType(ctx, sysType)
		if err != nil {
			return err
		}

		*(val.(*[]*model.Role)) = roles
		return nil
	})

	return resp, err
}

func (c *RoleCache) ClearRoles(ctx context.Context, sysType int64, userIds ...int64) error {
	keys := make([]string, len(userIds)+1)
	keys[0] = fmt.Sprintf("%s:%d", globals.ResourcesCacheKeyPrefix, sysType)
	for i, userId := range userIds {
		keys[i+1] = fmt.Sprintf("%s:%d:%d", globals.ResourcesCacheKeyPrefix, sysType, userId)
	}

	return c.cache.DelCtx(ctx, keys...)
}

func (c *RoleCache) ClearRolesBySysType(ctx context.Context, sysType int64) error {
	key := fmt.Sprintf("%s:%d", globals.RolesCacheKeyPrefix, sysType)
	return c.cache.DelCtx(ctx, key)
}

func (c *RoleCache) ClearRolesByUserIds(ctx context.Context, sysType int64, userIds ...int64) error {
	keys := slices.Map(userIds, func(userId int64) string {
		return fmt.Sprintf("%s:%d:%d", globals.RolesCacheKeyPrefix, sysType, userId)
	})
	return c.cache.DelCtx(ctx, keys...)
}
