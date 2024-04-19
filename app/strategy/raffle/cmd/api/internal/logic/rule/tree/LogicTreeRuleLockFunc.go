package tree

import (
	"context"
	"strconv"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

func RuleLockFunc(ctx context.Context, svc *svc.ServiceContext, userId string, strategyId int64, awardId int, ruleValue string) (TreeActionEntity, error) {
	//用户抽奖次数
	userCount := 100
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
				End:            false,
			},
		}, err
	}
	return TreeActionEntity{
		RuleLogicCheckTypeVO: TAKE_OVER,
		StrategyAwardVO: StrategyAwardVO{
			AwardId:        101,
			AwardRuleValue: "1,100",
			End:            true,
		},
	}, nil

}
