package logic

import (
	"context"
	"errors"
	"fmt"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/logic/rule/tree"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

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
	l.ctx = context.WithValue(l.ctx, "usedPoints", "2000")
	// 2.责任链接管 如果是黑名单，权重直接返回
	ChainFactory := chain.NewDefaultChainFactory(l.ctx, l.svcCtx)
	ChainFactory.OpenLogicChain(req.StrategyId)
	ChainStrategyAwardVO, err := ChainFactory.ExecLogicChain(req.StrategyId)
	if err != nil {
		logx.Error("ChainFactory.ExecLogicChain error:", err)
		return nil, err
	}
	if !chain.CheckStrategyAwardContinue(ChainStrategyAwardVO.LogicModel) {
		resp = &types.RaffleResponse{
			AwardId: ChainStrategyAwardVO.AwardId,
		}
		return
	}
	// 3.决策树
	// 3.1.获取决策
	strategyAward, err := l.svcCtx.StrategyAwardModel.QueryStrategyAward(l.ctx, req.StrategyId, ChainStrategyAwardVO.AwardId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return &types.RaffleResponse{
				AwardId: ChainStrategyAwardVO.AwardId,
			}, nil
		}
		logx.Error("QueryStrategyAward error:", err)
		return nil, err
	}
	// 3.2 获取ruleModels
	ruleModel := strategyAward.RuleModels.String
	if ruleModel == "" {
		return &types.RaffleResponse{AwardId: ChainStrategyAwardVO.AwardId}, nil
	}
	// 3.3 建树并且使用决策
	DefaultTreeFactor := tree.NewDefaultTreeFactor(l.ctx, l.svcCtx)
	enginee := DefaultTreeFactor.OpenLogicTree(strategyAward)
	TreeStrategyAwardVO, err := enginee.Process(l.ctx, user, req.StrategyId, ChainStrategyAwardVO.AwardId)
	if err != nil {
		logx.Error("NewDefaultTreeFactor.Process error:", err)
		return nil, err
	}
	// 返回
	resp = &types.RaffleResponse{
		AwardId: TreeStrategyAwardVO.AwardId,
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
