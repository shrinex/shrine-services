package repo

import (
	"core/member/rpc/internal/config"
	"core/member/rpc/internal/repo/cache"
	"core/member/rpc/internal/repo/db"
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
