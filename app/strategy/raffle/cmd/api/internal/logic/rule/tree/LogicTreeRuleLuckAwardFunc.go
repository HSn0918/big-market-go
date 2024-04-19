package tree

import (
	"context"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

func RuleLuckAwardFunc(ctx context.Context, svc *svc.ServiceContext, userId string, strategyId int64, awardId int, ruleValue string) (TreeActionEntity, error) {
	return TreeActionEntity{
		RuleLogicCheckTypeVO: TAKE_OVER,
		StrategyAwardVO: StrategyAwardVO{
			AwardId:        101,
			AwardRuleValue: "1,100",
			End:            false,
		},
	}, nil
}
