package chain

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/model"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

type DefaultChainFactory struct {
	loginChainGroup map[string]*LogicChain
	svcCtx          *svc.ServiceContext
	ctx             context.Context
}
type strategyAwardVO struct {
	AwardId    int64
	LogicModel string
}
type LogicChain struct {
	F    map[string]func(string, int64) strategyAwardVO
	Next *LogicChain
}

func NewDefaultChainFactory() DefaultChainFactory {
	return DefaultChainFactory{}
}
func (d DefaultChainFactory) OpenLogicChain(strategyId int64) *LogicChain {
	strategy, err := d.svcCtx.StrategyModel.QueryStrategy(d.ctx, strategyId)
	switch {
	case err == nil:
		break
	case errors.Is(err, model.ErrNotFound):
		// 填充默认规则
		return d.loginChainGroup[RULE_DEFAULT.Code()]

	default:
		logx.Error("QueryStrategy error:", err)
		return nil
	}
	ruleModels := strategy.GetRuleModels()
	loginChain, _ := d.loginChainGroup[ruleModels[0]]
	current := loginChain
	for i := 1; i < len(ruleModels); i++ {
		nextChain := d.loginChainGroup[ruleModels[i]]
		current.Next = nextChain
	}
	current.Next = d.loginChainGroup[RULE_DEFAULT.Code()]
	return loginChain
}
