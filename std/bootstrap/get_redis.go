package bootstrap

import (
	"fmt"
	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
	"io"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var clientManager = syncx.NewResourceManager()

func getRedis(cfg redis.RedisConf) (*red.Client, error) {
	if cfg.Type != redis.NodeType {
		return nil, fmt.Errorf("redis type '%s' is not supported", cfg.Type)
	}

	val, err := clientManager.GetResource(cfg.Host, func() (io.Closer, error) {
		return newRedis(cfg)
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}

func newRedis(cfg redis.RedisConf) (*red.Client, error) {
	return red.NewClient(&red.Options{
		Addr:         cfg.Host,
		Password:     cfg.Pass,
		DB:           defaultDatabase,
		MaxRetries:   maxRetries,
		MinIdleConns: idleConns,
	}), nil
}
