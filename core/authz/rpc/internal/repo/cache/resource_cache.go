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
	ResourceCache struct {
		roleCache *RoleCache
		cache     cache.Cache
		dao       model.ResourceModel
	}
)

func NewResourceCache(conf cache.CacheConf, roleCache *RoleCache, dao model.ResourceModel, opts ...cache.Option) *ResourceCache {
	cc := cache.New(conf, syncx.NewSingleFlight(), cache.NewStat("resource cache"), ErrNotFound, opts...)
	return &ResourceCache{dao: dao, cache: cc, roleCache: roleCache}
}

func (c *ResourceCache) ListResources(ctx context.Context, sysType, isAdmin, userId int64) ([]*model.Resource, error) {
	if isAdmin == globals.FlagTrue {
		return c.ListResourcesBySysType(ctx, sysType)
	}

	roles, err := c.roleCache.ListRoles(ctx, sysType, isAdmin, userId)
	if err != nil || slices.IsEmpty(roles) {
		return slices.Empty[*model.Resource](), err
	}

	// use map reduce to speed up
	resp, err := mr.MapReduce(func(source chan<- *model.Role) {
		for _, role := range roles {
			source <- role
		}
	}, func(item *model.Role, writer mr.Writer[[]*model.Resource], cancel func(error)) {
		res, err := c.ListResourcesByRoleId(ctx, sysType, item.RoleId)
		if err != nil {
			cancel(err)
		}
		writer.Write(res)
	}, func(pipe <-chan []*model.Resource, writer mr.Writer[[]*model.Resource], cancel func(error)) {
		res := slices.Empty[*model.Resource]()
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

func (c *ResourceCache) ListResourcesBySysType(ctx context.Context, sysType int64) ([]*model.Resource, error) {
	resp := slices.Empty[*model.Resource]()
	key := fmt.Sprintf("%s:%d", globals.ResourcesCacheKeyPrefix, sysType)
	err := c.cache.TakeCtx(ctx, &resp, key, func(val any) error {
		resources, err := c.dao.ListAllResourcesBySysType(ctx, sysType)
		if err != nil {
			return err
		}

		*(val.(*[]*model.Resource)) = resources
		return nil
	})

	return resp, err
}

func (c *ResourceCache) ListResourcesByRoleId(ctx context.Context, sysType, roleId int64) ([]*model.Resource, error) {
	resp := slices.Empty[*model.Resource]()
	key := fmt.Sprintf("%s:%d:%d", globals.ResourcesCacheKeyPrefix, sysType, roleId)
	err := c.cache.TakeCtx(ctx, &resp, key, func(val any) error {
		resources, err := c.dao.ListResourcesBySysTypeAndRoleId(ctx, sysType, roleId)
		if err != nil {
			return err
		}

		*(val.(*[]*model.Resource)) = resources
		return nil
	})

	return resp, err
}

func (c *ResourceCache) ClearResources(ctx context.Context, sysType int64, roleIds ...int64) error {
	keys := make([]string, len(roleIds)+1)
	keys[0] = fmt.Sprintf("%s:%d", globals.ResourcesCacheKeyPrefix, sysType)
	for i, roleId := range roleIds {
		keys[i+1] = fmt.Sprintf("%s:%d:%d", globals.ResourcesCacheKeyPrefix, sysType, roleId)
	}

	return c.cache.DelCtx(ctx, keys...)
}

func (c *ResourceCache) ClearResourcesBySysType(ctx context.Context, sysType int64) error {
	key := fmt.Sprintf("%s:%d", globals.ResourcesCacheKeyPrefix, sysType)
	return c.cache.DelCtx(ctx, key)
}

func (c *ResourceCache) ClearResourcesByRoleIds(ctx context.Context, sysType int64, roleIds ...int64) error {
	keys := slices.Map(roleIds, func(roleId int64) string {
		return fmt.Sprintf("%s:%d:%d", globals.ResourcesCacheKeyPrefix, sysType, roleId)
	})
	return c.cache.DelCtx(ctx, keys...)
}
