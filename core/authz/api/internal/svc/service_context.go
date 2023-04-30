package svc

import (
	"core/authz/api/internal/config"
	"core/authz/api/internal/realms"
	"core/authz/rpc/service"
	"github.com/shrinex/shield/authc"
	"github.com/shrinex/shield/authz"
	"github.com/shrinex/shield/codec"
	"github.com/shrinex/shield/security"
	"github.com/zeromicro/go-zero/zrpc"
	"shrine/std/authx"
)

type ServiceContext struct {
	Config   config.Config
	Subject  security.Subject
	AuthzRpc service.Service
}

func NewServiceContext(cfg config.Config) *ServiceContext {
	authzRpc := service.NewService(zrpc.MustNewClient(cfg.AuthzRpc))
	repository := authx.NewRepository(cfg.Redis, authx.FlushModeOnSave, codec.JSON,
		security.GetGlobalOptions().GetTimeout(), security.GetGlobalOptions().GetIdleTimeout())
	return &ServiceContext{
		Config:   cfg,
		AuthzRpc: authzRpc,
		Subject: security.NewBuilder[*authx.RedisSession]().
			Authenticator(authc.NewAuthenticator(
				authx.NewBearerAuthRealm(repository))).
			Authorizer(authz.NewAuthorizer(
				realms.NewAuthzRpcRealm(authzRpc))).
			Registry(authx.NewRegistry(repository)).
			Repository(repository).
			Build(),
	}
}
