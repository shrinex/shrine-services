package main

import (
	"core/authz/api/internal/realms"
	"flag"
	"fmt"
	"github.com/shrinex/shield-web/chain"
	"github.com/shrinex/shield/authz"
	"github.com/shrinex/shield/security"
	"net/http"

	"core/authz/api/internal/config"
	"core/authz/api/internal/handler"
	"core/authz/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	_ "shrine/std/valid" // register global validator
)

var configFile = flag.String("f", "etc/main.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	server.Use(chain.NewBuilder().
		Subject().
		Use(ctx.Subject).
		And().
		SessionManagement().
		And().
		BearerAuth().
		AnyRequests().
		And().
		AuthorizeRequests().
		AnyRequests().
		Authenticated().
		And().
		AuthorizeRequests().
		AnyRequests().
		HasAuthorityFunc(func(r *http.Request, _ security.Subject) authz.Authority {
			return realms.NewAuthzAuthority(r.Method, r.URL.Path)
		}).
		And().
		Build(),
	)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
