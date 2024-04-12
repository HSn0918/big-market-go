package logic

import (
	"context"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RaffleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRaffleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RaffleLogic {
	return &RaffleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RaffleLogic) Raffle(req *types.RaffleRequest) (resp *types.RaffleResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
