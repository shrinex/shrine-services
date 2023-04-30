package repo

import (
	"core/authc/rpc/internal/config"
	"core/authc/rpc/internal/repo/cache"
	"core/authc/rpc/internal/repo/db"
)

type Aggregate struct {
	DB    *db.Repository
	Cache *cache.Repository
}

func NewAggregate(cfg config.Config) *Aggregate {
	return &Aggregate{
		DB:    db.NewRepository(cfg),
		Cache: cache.NewRepository(cfg),
	}
}
