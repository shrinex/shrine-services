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

func (c *RoleCache) ListRoles(ctx context.Context, sysType, userId, isAdmin int64) (roles []*model.Role, err error) {
	resp := slices.Empty[*model.Role]()
	key := fmt.Sprintf("%s:%d:%d", globals.RolesCacheKeyPrefix, sysType, userId)
	err = c.cache.TakeCtx(ctx, &resp, key, func(val any) error {
		if isAdmin == globals.FlagTrue {
			roles, err = c.dao.ListAllRolesBySysType(ctx, sysType)
		} else {
			roles, err = c.dao.ListRolesByUserIdAndSysType(ctx, userId, sysType)
		}

		if err != nil {
			return err
		}

		*(val.(*[]*model.Role)) = roles
		return nil
	})

	return resp, err
}

func (c *RoleCache) ClearRoles(ctx context.Context, sysType, userId int64) error {
	key := fmt.Sprintf("%s:%d:%d", globals.RolesCacheKeyPrefix, sysType, userId)
	return c.cache.DelCtx(ctx, key)
}
