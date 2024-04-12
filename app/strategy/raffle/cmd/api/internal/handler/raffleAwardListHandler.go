package handler

import (
	"net/http"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/logic"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RaffleAwardListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RaffleAwardListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRaffleAwardListLogic(r.Context(), svcCtx)
		resp, err := l.RaffleAwardList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
