package repo

import (
	"unit/product/rpc/internal/config"
	"unit/product/rpc/internal/repo/cache"
	"unit/product/rpc/internal/repo/db"
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
