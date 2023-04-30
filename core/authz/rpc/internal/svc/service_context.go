package svc

import (
	"core/authz/rpc/internal/config"
	"core/authz/rpc/internal/repo"
	"shrine/std/leaf"
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
		Leaf:      leaf.NewSnowflake(leaf.Authz),
	}
}
