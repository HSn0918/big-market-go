package logic

import (
	"context"
	"errors"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/model"
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
	// 1.查询策略配置
	var StrategyAwardList []*model.StrategyAward
	StrategyAwardList, err = l.svcCtx.StrategyAwardModel.QueryStrategyAwardList(l.ctx, req.StrategyId)
	// 2 缓存奖品库存【用于decr扣减库存使用】
	for _, v := range StrategyAwardList {
		err = l.cacheStrategyAwardCount(req.StrategyId, int(v.AwardId), int(v.AwardCount))
		if err != nil {
			logx.Error("cacheStrategyAwardCount error:", err)
		}
	}
	// 3.1 默认装配配置【全量抽奖概率】
	l.assembleLotteryStrategy(req.StrategyId, StrategyAwardList)
	resp = &types.StrategyArmoryResponse{IsSuccess: true}
	// 3.2 装配权重策略配置
	strategy, err := l.svcCtx.StrategyModel.QueryStrategy(l.ctx, req.StrategyId)
	if err != nil {
		// 3.2.1 "rule_weight"规则不存在，则直接返回
		if errors.Is(err, model.ErrNotFound) {
			return resp, nil
		}
		resp = &types.StrategyArmoryResponse{IsSuccess: false}
		logx.Error("QueryStrategy error:", err)
		return resp, err
	}
	// 3.2.2 "rule_weight"规则存在，则需要进行权重策略配置

	strategyRule, err := l.svcCtx.StrategyRuleModel.QueryStrategyRule(l.ctx, req.StrategyId, strategy.GetRuleWeight())
	// 业务异常，策略规则中 rule_weight 权重规则已适用但未配置
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return resp, nil
		}
		logx.Error("QueryStrategyRule error:", err)
		return resp, err
	}
	// 3.2.3 装配策略奖品概率查找表
	RuleWeightValueMap := strategyRule.GetRuleWeightValues()
	for k, v := range RuleWeightValueMap {
		ruleWeightValues := v
		StrategyAwardListFilter := cloneAndFilterStrategyAward(StrategyAwardList, ruleWeightValues)
		_, err = l.assembleLotteryStrategyByRuleWeight(req.StrategyId, k, StrategyAwardListFilter)
		if err != nil {
			resp = &types.StrategyArmoryResponse{IsSuccess: false}
			logx.Error("assembleLotteryStrategyByRuleWeight error:", err)
			return resp, err
		}
	}

	return resp, nil
}
