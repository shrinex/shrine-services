package cache

import (
	"context"
	"core/authz/proto/model"
	"fmt"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"shrine/std/globals"
	"shrine/std/utils/slices"
)

type (
	MenuCache struct {
		roleCache *RoleCache
		cache     cache.Cache
		dao       model.MenuModel
	}
)

func NewMenuCache(conf cache.CacheConf, roleCache *RoleCache, dao model.MenuModel, opts ...cache.Option) *MenuCache {
	cc := cache.New(conf, syncx.NewSingleFlight(), cache.NewStat("menu cache"), ErrNotFound, opts...)
	return &MenuCache{dao: dao, cache: cc, roleCache: roleCache}
}

func (c *MenuCache) ListMenus(ctx context.Context, sysType, isAdmin, userId int64) (menus []*model.Menu, err error) {
	if isAdmin == globals.FlagTrue {
		return c.ListMenusBySysType(ctx, sysType)
	}

	roles, err := c.roleCache.ListRoles(ctx, sysType, isAdmin, userId)
	if err != nil || slices.IsEmpty(roles) {
		return slices.Empty[*model.Menu](), err
	}

	// use map reduce to speed up
	resp, err := mr.MapReduce(func(source chan<- *model.Role) {
		for _, role := range roles {
			source <- role
		}
	}, func(item *model.Role, writer mr.Writer[[]*model.Menu], cancel func(error)) {
		res, err := c.ListMenusByRoleId(ctx, sysType, item.RoleId)
		if err != nil {
			cancel(err)
		}
		writer.Write(res)
	}, func(pipe <-chan []*model.Menu, writer mr.Writer[[]*model.Menu], cancel func(error)) {
		res := slices.Empty[*model.Menu]()
		for item := range pipe {
			res = append(res, item...)
		}
		writer.Write(res)
	}, mr.WithWorkers(mathx.MinInt(16, len(roles))))
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (c *MenuCache) ListMenusBySysType(ctx context.Context, sysType int64) ([]*model.Menu, error) {
	resp := slices.Empty[*model.Menu]()
	key := fmt.Sprintf("%s:%d", globals.MenusCacheKeyPrefix, sysType)
	err := c.cache.TakeCtx(ctx, &resp, key, func(val any) error {
		menus, err := c.dao.ListAllMenusBySysType(ctx, sysType)
		if err != nil {
			return err
		}

		*(val.(*[]*model.Menu)) = menus
		return nil
	})

	return resp, err
}

func (c *MenuCache) ListMenusByRoleId(ctx context.Context, sysType, roleId int64) ([]*model.Menu, error) {
	resp := slices.Empty[*model.Menu]()
	key := fmt.Sprintf("%s:%d:%d", globals.MenusCacheKeyPrefix, sysType, roleId)
	err := c.cache.TakeCtx(ctx, &resp, key, func(val any) error {
		menus, err := c.dao.ListMenusBySysTypeAndRoleId(ctx, sysType, roleId)
		if err != nil {
			return err
		}

		*(val.(*[]*model.Menu)) = menus
		return nil
	})

	return resp, err
}

func (c *MenuCache) ClearMenus(ctx context.Context, sysType int64, roleIds ...int64) error {
	keys := make([]string, len(roleIds)+1)
	keys[0] = fmt.Sprintf("%s:%d", globals.MenusCacheKeyPrefix, sysType)
	for i, roleId := range roleIds {
		keys[i+1] = fmt.Sprintf("%s:%d:%d", globals.MenusCacheKeyPrefix, sysType, roleId)
	}

	return c.cache.DelCtx(ctx, keys...)
}

func (c *MenuCache) ClearMenusBySysType(ctx context.Context, sysType int64) error {
	key := fmt.Sprintf("%s:%d", globals.MenusCacheKeyPrefix, sysType)
	return c.cache.DelCtx(ctx, key)
}

func (c *MenuCache) ClearMenusByRoleIds(ctx context.Context, sysType int64, roleIds ...int64) error {
	keys := slices.Map(roleIds, func(roleId int64) string {
		return fmt.Sprintf("%s:%d:%d", globals.MenusCacheKeyPrefix, sysType, roleId)
	})
	return c.cache.DelCtx(ctx, keys...)
}
