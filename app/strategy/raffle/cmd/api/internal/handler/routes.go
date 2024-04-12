// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/query_raffle_award_list",
				Handler: RaffleAwardListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/random_raffle",
				Handler: RaffleHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/strategy_armory",
				Handler: StrategyArmoryHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1/raffle"),
	)
}
