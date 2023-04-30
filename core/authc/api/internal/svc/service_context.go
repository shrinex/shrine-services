package svc

import (
	"core/authc/api/internal/config"
	"core/authc/api/internal/realms"
	"core/authc/rpc/service"
	"github.com/shrinex/shield/authc"
	"github.com/shrinex/shield/authz"
	"github.com/shrinex/shield/codec"
	"github.com/shrinex/shield/security"
	"github.com/zeromicro/go-zero/zrpc"
	"shrine/std/authx"
)

type ServiceContext struct {
	Config   config.Config
	AuthcRpc service.Service
	Subject  security.Subject
}

func NewServiceContext(cfg config.Config) *ServiceContext {
	authcRpc := service.NewService(zrpc.MustNewClient(cfg.AuthcRpc))
	repository := authx.NewRepository(cfg.Redis, authx.FlushModeOnSave, codec.JSON,
		security.GetGlobalOptions().GetTimeout(), security.GetGlobalOptions().GetIdleTimeout())
	return &ServiceContext{
		Config:   cfg,
		AuthcRpc: authcRpc,
		Subject: security.NewBuilder[*authx.RedisSession]().
			Authenticator(authc.NewAuthenticator(
				authx.NewBearerAuthRealm(repository),
				realms.NewAuthcRpcRealm(authcRpc))).
			Authorizer(authz.NoopAuthorizer).
			Registry(authx.NewRegistry(repository)).
			Repository(repository).
			Build(),
	}
}
