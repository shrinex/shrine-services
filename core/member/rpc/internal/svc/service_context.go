package svc

import (
	"core/member/rpc/internal/config"
	"core/member/rpc/internal/repo"
	"shrine/std/leaf"
)

type ServiceContext struct {
	Config config.Config
	Leaf   leaf.Snowflake
	*repo.Aggregate
}

func NewServiceContext(cfg config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    cfg,
		Leaf:      leaf.NewSnowflake(leaf.Authc),
		Aggregate: repo.NewAggregate(cfg),
	}
}
