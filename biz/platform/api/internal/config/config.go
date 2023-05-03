package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"shrine/std/conf/dtm"
	"shrine/std/conf/rdb"
)

type Config struct {
	rest.RestConf
	Redis    redis.RedisConf
	MySQL    rdb.MySQLConf
	AuthcRpc zrpc.RpcClientConf
	AuthzRpc zrpc.RpcClientConf
	ShopRpc  zrpc.RpcClientConf
	Dtm      dtm.DtmConf
}
