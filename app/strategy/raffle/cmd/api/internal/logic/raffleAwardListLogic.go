package logic

import (
	"context"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RaffleAwardListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRaffleAwardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RaffleAwardListLogic {
	return &RaffleAwardListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RaffleAwardListLogic) RaffleAwardList(req *types.RaffleAwardListRequest) (resp *types.RaffleAwardListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
