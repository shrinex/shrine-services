package main

import (
	"flag"
	"fmt"
	"github.com/shrinex/shield-web/chain"

	"core/authc/api/internal/config"
	"core/authc/api/internal/handler"
	"core/authc/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
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
		AntExcludes("/**/login").
		And().
		AuthorizeRequests().
		AnyRequests().
		AntExcludes("/**/login").
		Authenticated().
		And().
		Build(),
	)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
