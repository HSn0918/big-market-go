package tree

import "context"

type StrategyAwardVO struct {
	AwardId        int
	AwardRuleValue string
	End            bool
}

type RuleLogicCheckTypeVO string
type LogicModel string
type TreeActionEntity struct {
	RuleLogicCheckTypeVO RuleLogicCheckTypeVO
	StrategyAwardVO      StrategyAwardVO
}
type Logic func(ctx context.Context, userId string, strategyId int64, awardId int, ruleValue string) (TreeActionEntity, error)

const (
	RULE_LOCK       LogicModel           = "限定用户已完成N次抽奖后解锁"
	RULE_STOCK      LogicModel           = "库存扣减规则"
	RULE_LUCK_AWARD LogicModel           = "抽奖概率规则"
	ALLOW           RuleLogicCheckTypeVO = "放行；执行后续的流程，不受规则引擎影响"
	TAKE_OVER       RuleLogicCheckTypeVO = "接管；后续的流程，受规则引擎执行结果影响"
	EQUAL           string               = "EQUAL"
	GT              string               = "GT"
	LT              string               = "LT"
	GE              string               = "GE"
	LE              string               = "LE"
)

func (m RuleLogicCheckTypeVO) Code() string {
	switch m {
	case ALLOW:
		return "0000"
	case TAKE_OVER:
		return "0001"
	default:
		return "unknown"
	}

}
func (m LogicModel) Code() string {
	switch m {
	case RULE_LOCK:
		return "rule_lock"
	case RULE_STOCK:
		return "rule_stock"
	case RULE_LUCK_AWARD:
		return "rule_luck_award"
	default:
		return "unknown"
	}
}

type treeActionEntity struct {
	RuleLogicCheckTypeVO RuleLogicCheckTypeVO
	StrategyAwardVO      StrategyAwardVO
}
