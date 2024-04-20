package logic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/code"

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
	strategyAwardList, err := l.svcCtx.StrategyAwardModel.QueryStrategyAwardList(l.ctx, req.StrategyId)
	if err != nil {
		// 如果查询过程中发生错误，根据错误类型返回相应的错误码
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, code.RaffleAwardNotFound // 假设 code.RaffleAwardNotFound 是定义的一个错误码
		}
		return nil, err // 返回其他类型的错误
	}
	// 如果查询结果为空，返回空列表错误码
	if len(strategyAwardList) == 0 {
		return nil, code.RaffleAwardEmpty
	}
	// 准备响应数据
	raffleAwardList := make([]*types.RaffleAward, 0, len(strategyAwardList))
	for _, award := range strategyAwardList {
		// 将数据库模型转换为响应模型
		raffleAwardList = append(raffleAwardList, &types.RaffleAward{
			AwardId:       int(award.AwardId),
			AwardTitle:    award.AwardTitle,
			AwardSubtitle: award.AwardSubtitle.String, // 假设前端需要的是 int 类型
			Sort:          int(award.Sort),
		})
	}
	// 创建响应体并返回
	resp = &types.RaffleAwardListResponse{
		RaffleAwardList: raffleAwardList,
	}
	return resp, nil
}
