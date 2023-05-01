package repo

import (
	"unit/shop/rpc/internal/config"
	"unit/shop/rpc/internal/repo/cache"
	"unit/shop/rpc/internal/repo/db"
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
