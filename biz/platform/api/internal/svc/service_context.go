package svc

import (
	"biz/platform/api/internal/config"
	authcpkg "core/authc/rpc/service"
	"core/authz/rpc/realms"
	authzpkg "core/authz/rpc/service"
	"github.com/shrinex/shield/authc"
	"github.com/shrinex/shield/authz"
	"github.com/shrinex/shield/codec"
	"github.com/shrinex/shield/security"
	"github.com/zeromicro/go-zero/zrpc"
	"shrine/std/authx"
	"shrine/std/leaf"
	shoppkg "unit/shop/rpc/service"
)

type ServiceContext struct {
	Config   config.Config
	Subject  security.Subject
	AuthcRpc authcpkg.Service
	AuthzRpc authzpkg.Service
	ShopRpc  shoppkg.Service
	Leaf     leaf.Snowflake
}

func NewServiceContext(cfg config.Config) *ServiceContext {
	authzRpc := authzpkg.NewService(zrpc.MustNewClient(cfg.AuthzRpc))
	repository := authx.NewRepository(cfg.Redis, authx.FlushModeOnSave, codec.JSON,
		security.GetGlobalOptions().GetTimeout(), security.GetGlobalOptions().GetIdleTimeout())
	return &ServiceContext{
		Config:   cfg,
		AuthzRpc: authzRpc,
		Leaf:     leaf.NewSnowflake(leaf.Dtm),
		AuthcRpc: authcpkg.NewService(zrpc.MustNewClient(cfg.AuthcRpc)),
		ShopRpc:  shoppkg.NewService(zrpc.MustNewClient(cfg.ShopRpc)),
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
