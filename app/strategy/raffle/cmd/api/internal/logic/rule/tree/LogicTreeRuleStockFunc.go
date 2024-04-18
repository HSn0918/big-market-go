package tree

import "context"

func RuleStockFunc(ctx context.Context, userId string, strategyId int64, awardId int, ruleValue string) (TreeActionEntity, error) {
	return TreeActionEntity{
		RuleLogicCheckTypeVO: TAKE_OVER,
		StrategyAwardVO: StrategyAwardVO{
			AwardId:        awardId,
			AwardRuleValue: ruleValue,
			End:            false,
		},
	}, nil
}
