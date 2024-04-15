package chain

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/model"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

type strategyAwardVO struct {
	AwardId    int64
	LogicModel string
}
type LogicChainFunc func(string, int64) strategyAwardVO
type DefaultChainFactory struct {
	loginChainGroup map[string]LogicChainFunc
	Chain           *LogicChain
	svcCtx          *svc.ServiceContext
	ctx             context.Context
}

type LogicChain struct {
	Func LogicChainFunc
	Next *LogicChain
}

func NewDefaultChainFactory(ctx context.Context, svcCtx *svc.ServiceContext) DefaultChainFactory {
	return DefaultChainFactory{
		ctx:             ctx,
		svcCtx:          svcCtx,
		loginChainGroup: NewLogicChainGroup(),
		Chain:           new(LogicChain),
	}
}
func NewLogicChainGroup() map[string]LogicChainFunc {
	mp := make(map[string]LogicChainFunc, 0)
	// todo:func 实现
	mp[RULE_DEFAULT.Code()] = nil
	mp[RULE_BLACKLIST.Code()] = nil
	mp[RULE_WEIGHT.Code()] = nil
	return mp
}
func (d *DefaultChainFactory) OpenLogicChain(strategyId int64) *LogicChain {
	strategy, err := d.svcCtx.StrategyModel.QueryStrategy(d.ctx, strategyId)
	switch {
	case err == nil:
		break
	case errors.Is(err, model.ErrNotFound):
		// 填充默认规则
		d.Chain.Func = d.loginChainGroup[RULE_DEFAULT.Code()]
		return d.Chain
	default:
		logx.Error("QueryStrategy error:", err)
		return nil
	}
	ruleModels := strategy.GetRuleModels()
	logicChain := new(LogicChain)
	current := logicChain
	for i, ruleModel := range ruleModels {
		funcNode, exists := d.loginChainGroup[ruleModel]
		if !exists {
			continue // 如果规则模型不存在，跳过当前迭代
		}
		current.Func = funcNode
		if i < len(ruleModels)-1 {
			current.Next = new(LogicChain) // 为下一个规则模型创建新节点
			current = current.Next
		}

	}
	// 设置最后一个逻辑链节点的处理函数
	current.Next = new(LogicChain)
	current.Next.Func = d.loginChainGroup[RULE_DEFAULT.Code()]
	d.Chain = logicChain
	return d.Chain
}
