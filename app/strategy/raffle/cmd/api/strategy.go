package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/hsn0918/BigMarket/pkg/xcode"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/logic"
	"github.com/robfig/cron/v3"

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

	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()
	ctx := svc.NewServiceContext(c)

	httpx.SetErrorHandler(xcode.ErrHandler)
	httpx.SetOkHandler(xcode.OkHandler)
	handler.RegisterHandlers(server, ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	// 定时任务
	cronJob := logic.NewCronJob(context.Background(), ctx, cron.New(cron.WithSeconds()))
	cronJob.Job.Start()
	server.Start()
}
