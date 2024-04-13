package main

import (
	"flag"
	"fmt"

	"github.com/hsn0918/BigMarket/pkg/xcode"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/config"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/handler"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/strategy-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	httpx.SetErrorHandler(xcode.ErrHandler)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
