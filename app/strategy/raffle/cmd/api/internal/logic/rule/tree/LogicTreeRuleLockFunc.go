package tree

import (
	"context"
	"strconv"
)

func RuleLockFunc(ctx context.Context, userId string, strategyId int64, awardId int, ruleValue string) (TreeActionEntity, error) {
	//userCount := ctx.Value("userCount").(int)
	userCount := 0
	ruleValueInt, err := strconv.Atoi(ruleValue)
	if err != nil {
		return TreeActionEntity{}, err
	}
	if userCount >= ruleValueInt {
		return TreeActionEntity{
			RuleLogicCheckTypeVO: ALLOW,
			StrategyAwardVO: StrategyAwardVO{
				AwardId:        awardId,
				AwardRuleValue: ruleValue,
				End:            true,
			},
		}, err
	}
	return TreeActionEntity{
		RuleLogicCheckTypeVO: TAKE_OVER,
		StrategyAwardVO: StrategyAwardVO{
			AwardId:        100,
			AwardRuleValue: ruleValue,
			End:            false,
		},
	}, nil

}
