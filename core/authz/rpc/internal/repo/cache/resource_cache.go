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

func (c *ResourceCache) ListResources(ctx context.Context, sysType, userId, isAdmin int64) ([]*model.Resource, error) {
	if isAdmin == globals.FlagTrue {
		return c.ListResourcesBySysType(ctx, sysType)
	}

	roles, err := c.roleCache.ListRoles(ctx, sysType, userId, isAdmin)
	if err != nil || slices.IsEmpty(roles) {
		return slices.Empty[*model.Resource](), err
	}

	// use map reduce to speed up
	resp, err := mr.MapReduce(func(source chan<- *model.Role) {
		for _, role := range roles {
			source <- role
		}
	}, func(item *model.Role, writer mr.Writer[[]*model.Resource], cancel func(error)) {
		res, err := c.ListResourcesByRole(ctx, item)
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

func (c *ResourceCache) ListResourcesByRole(ctx context.Context, role *model.Role) ([]*model.Resource, error) {
	return c.ListResourcesByRoleId(ctx, role.RoleId)
}

func (c *ResourceCache) ListResourcesByRoleId(ctx context.Context, roleId int64) ([]*model.Resource, error) {
	resp := slices.Empty[*model.Resource]()
	key := fmt.Sprintf("%s:%d", globals.ResourcesCacheKeyPrefix, roleId)
	err := c.cache.TakeCtx(ctx, &resp, key, func(val any) error {
		resources, err := c.dao.ListResourcesByRoleId(ctx, roleId)
		if err != nil {
			return err
		}

		*(val.(*[]*model.Resource)) = resources
		return nil
	})

	return resp, err
}

func (c *ResourceCache) ClearResources(ctx context.Context, sysType, userId, isAdmin int64) error {
	if isAdmin == globals.FlagTrue {
		return c.ClearResourcesBySysType(ctx, sysType)
	}

	roles, err := c.roleCache.ListRoles(ctx, sysType, userId, isAdmin)
	if err != nil {
		return err
	}

	var keys []string
	for _, role := range roles {
		key := fmt.Sprintf("%s:%d", globals.ResourcesCacheKeyPrefix, role.RoleId)
		keys = append(keys, key)
	}

	return c.cache.DelCtx(ctx, keys...)
}

func (c *ResourceCache) ClearResourcesBySysType(ctx context.Context, sysType int64) error {
	key := fmt.Sprintf("%s:%d", globals.ResourcesCacheKeyPrefix, sysType)
	return c.cache.DelCtx(ctx, key)
}

func (c *ResourceCache) ClearResourcesByRole(ctx context.Context, role *model.Role) error {
	return c.ClearResourcesByRoleId(ctx, role.RoleId)
}

func (c *ResourceCache) ClearResourcesByRoleId(ctx context.Context, roleId int64) error {
	key := fmt.Sprintf("%s:%d", globals.ResourcesCacheKeyPrefix, roleId)
	return c.cache.DelCtx(ctx, key)
}
