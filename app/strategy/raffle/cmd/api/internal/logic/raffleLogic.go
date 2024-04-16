package logic

import (
	"context"
	"fmt"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/logic/rule/chain"

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
	// 1. 参数校验
	strategyId := req.StrategyId
	// userId 从context中获取
	// 这里假设userId为"system"
	user := "user00"
	if user == "" || strategyId <= 0 {
		return nil, fmt.Errorf("invalid params")
	}
	l.ctx = context.WithValue(l.ctx, "user", user)
	l.ctx = context.WithValue(l.ctx, "usedPoints", "100")

	//todo 前置规则 责任链
	ChainFactory := chain.NewDefaultChainFactory(l.ctx, l.svcCtx)
	ChainFactory.OpenLogicChain(req.StrategyId)
	awardId, err := ChainFactory.ExecLogicGroup(req.StrategyId, user)
	// 返回
	resp = &types.RaffleResponse{
		AwardId: awardId,
	}
	return
}

//func (l *RaffleLogic) getRandomAwardId(StrategyId int64) (awardId int, err error) {
//	// 1.从redis中取RateRange
//	cacheRateRange := fmt.Sprintf(redis.StrategyRateRangeSize, StrategyId)
//	rateRangeStr, err := l.svcCtx.BizRedis.Get(cacheRateRange)
//	if err != nil {
//		return -1, err
//	}
//	rateRange, err := strconv.Atoi(rateRangeStr)
//	if err != nil {
//		return -1, err
//	}
//	randInt := rand.IntN(rateRange)
//	// 2.从redis中取AwardId
//
//	cacheStrategy := fmt.Sprintf(redis.StrategyRateRange, StrategyId)
//
//	awardIdStr, err := l.svcCtx.BizRedis.HgetCtx(l.ctx, cacheStrategy, strconv.Itoa(randInt))
//	if err != nil {
//		return -1, err
//	}
//	awardId, err = strconv.Atoi(awardIdStr)
//	if err != nil {
//		return -1, err
//	}
//	return
//}
