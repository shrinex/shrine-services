package svc

import (
	"shrine/std/leaf"
	"unit/shop/rpc/internal/config"
	"unit/shop/rpc/internal/repo"
)

type ServiceContext struct {
	Config config.Config
	*repo.Aggregate
	Leaf leaf.Snowflake
}

func NewServiceContext(cfg config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    cfg,
		Aggregate: repo.NewAggregate(cfg),
		Leaf:      leaf.NewSnowflake(leaf.Merchant),
	}
}
