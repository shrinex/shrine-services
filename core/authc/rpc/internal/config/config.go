package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"shrine/std/conf/rdb"
)

type Config struct {
	zrpc.RpcServerConf
	Cache cache.CacheConf
	MySQL rdb.MySQL
}
