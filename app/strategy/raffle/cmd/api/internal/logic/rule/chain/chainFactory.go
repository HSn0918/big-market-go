package chain

import (
	"context"
	"errors"
	"fmt"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/model"

	"math/rand/v2"
	"strconv"

	"github.com/hsn0918/BigMarket/common"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

type DefaultChainFactory struct {
	loginChainGroup map[string]LogicChainFunc
	Chain           *LogicChain
	ctx             context.Context
	svcCtx          *svc.ServiceContext
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
	mp := make(map[string]LogicChainFunc)
	mp[RULE_DEFAULT.Code()] = DefaultFunc
	mp[RULE_BLACKLIST.Code()] = BlackFunc
	mp[RULE_WEIGHT.Code()] = WeightFunc
	return mp
}
func (d *DefaultChainFactory) OpenLogicChain(strategyId int64) *LogicChain {
	strategy, err := d.svcCtx.StrategyModel.QueryStrategy(d.ctx, strategyId)
	if err != nil {
		logx.Error("QueryStrategy error:", err)
		if errors.Is(err, model.ErrNotFound) {
			return &LogicChain{Func: d.loginChainGroup[RULE_DEFAULT.Code()]}
		}
		return nil
	}
	ruleModels := strategy.GetRuleModels()
	var logicChain *LogicChain
	var current *LogicChain
	for _, ruleModel := range ruleModels {
		funcNode, exists := d.loginChainGroup[ruleModel]
		if !exists {
			continue // 忽略不存在的规则
		}
		if logicChain == nil {
			logicChain = &LogicChain{Func: funcNode}
			current = logicChain
		} else {
			current.Next = &LogicChain{Func: funcNode}
			current = current.Next
		}
	}
	if current == nil {
		return &LogicChain{Func: d.loginChainGroup[RULE_DEFAULT.Code()]}
	}
	current.Next = &LogicChain{Func: d.loginChainGroup[RULE_DEFAULT.Code()]}
	d.Chain = logicChain
	return logicChain
}
func (d *DefaultChainFactory) ExecLogicChain(strategyId int64) (strategyAwardVO StrategyAwardVO, err error) {
	current := d.Chain
	for current != nil {
		if current.Func != nil {
			strategyAwardVO, err = current.Func(d.ctx, d.svcCtx, strategyId)
			if err != nil {
				logx.Error("func error:", err)
				return StrategyAwardVO{
					AwardId:    0,
					LogicModel: RULE_ERROR.Code(),
					End:        false,
				}, err
			}
			if strategyAwardVO.End {
				return strategyAwardVO, nil
			}
		}
		current = current.Next
	}
	return strategyAwardVO, err
}
func (d *DefaultChainFactory) getRandomAwardId(StrategyId int64) (awardId int, err error) {
	// 1.从redis中取RateRange
	cacheRateRange := fmt.Sprintf(common.StrategyRateRangeSize, StrategyId)
	rateRangeStr, err := d.svcCtx.BizRedis.Get(cacheRateRange)
	if err != nil {
		return -1, err
	}
	rateRange, err := strconv.Atoi(rateRangeStr)
	if err != nil {
		return -1, err
	}
	randInt := rand.IntN(rateRange)
	// 2.从redis中取AwardId
	cacheStrategy := fmt.Sprintf(common.StrategyRateRange, StrategyId)
	awardIdStr, err := d.svcCtx.BizRedis.HgetCtx(d.ctx, cacheStrategy, strconv.Itoa(randInt))
	if err != nil {
		return -1, err
	}
	awardId, err = strconv.Atoi(awardIdStr)
	if err != nil {
		return -1, err
	}
	return
}
