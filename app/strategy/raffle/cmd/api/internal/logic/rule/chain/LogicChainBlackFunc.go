package chain

import (
	"context"
	"strconv"
	"strings"

	"github.com/hsn0918/BigMarket/common"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

func BlackFunc(ctx context.Context, svc *svc.ServiceContext, strategyId int64) (StrategyAwardVO, error) {
	user := ctx.Value("user").(string)
	// 1.查询规则值配置
	strategy, err := svc.StrategyRuleModel.QueryStrategyRule(ctx, strategyId, RULE_BLACKLIST.Code())
	if err != nil {
		return StrategyAwardVO{}, err
	}
	ruleValue := strategy.RuleValue
	ruleValueGroup := strings.Split(ruleValue, common.COLON)
	awardId, _ := strconv.Atoi(ruleValueGroup[0])
	blackUserGroup := ruleValueGroup[1]
	blackUsers := strings.Split(blackUserGroup, common.SPLIT)
	// 2.判断是否在黑名单
	for _, blackUser := range blackUsers {
		if user == blackUser {
			return StrategyAwardVO{
				AwardId:    awardId,
				End:        true,
				LogicModel: RULE_BLACKLIST.Code(),
			}, nil
		}
	}
	// 3.不在黑名单中
	return StrategyAwardVO{
		AwardId:    101,
		End:        false,
		LogicModel: RULE_BLACKLIST.Code(),
	}, nil
}
