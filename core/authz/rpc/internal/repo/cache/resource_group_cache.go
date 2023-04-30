package cache

import (
	"core/authz/proto/model"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
)

type (
	ResourceGroupCache struct {
		cache cache.Cache
		dao   model.ResourceGroupModel
	}
)

func NewResourceGroupCache(conf cache.CacheConf, dao model.ResourceGroupModel, opts ...cache.Option) *ResourceGroupCache {
	cc := cache.New(conf, syncx.NewSingleFlight(), cache.NewStat("resource-group cache"), ErrNotFound, opts...)
	return &ResourceGroupCache{dao: dao, cache: cc}
}
