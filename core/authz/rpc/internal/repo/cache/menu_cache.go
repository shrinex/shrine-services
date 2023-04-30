package cache

import (
	"core/authz/proto/model"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
)

type (
	MenuCache struct {
		cache cache.Cache
		dao   model.MenuModel
	}
)

func NewMenuCache(conf cache.CacheConf, dao model.MenuModel, opts ...cache.Option) *MenuCache {
	cc := cache.New(conf, syncx.NewSingleFlight(), cache.NewStat("menu cache"), ErrNotFound, opts...)
	return &MenuCache{dao: dao, cache: cc}
}
