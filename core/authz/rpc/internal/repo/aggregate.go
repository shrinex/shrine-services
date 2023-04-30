package repo

import (
	"core/authz/rpc/internal/config"
	"core/authz/rpc/internal/repo/cache"
	"core/authz/rpc/internal/repo/db"
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
