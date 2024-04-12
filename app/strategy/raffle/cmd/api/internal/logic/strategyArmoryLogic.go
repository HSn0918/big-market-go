package logic

import (
	"context"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StrategyArmoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyArmoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyArmoryLogic {
	return &StrategyArmoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyArmoryLogic) StrategyArmory(req *types.StrategyArmoryRequest) (resp *types.StrategyArmoryResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
