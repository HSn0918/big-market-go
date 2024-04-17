package tree

type StrategyAwardVO struct {
	AwardId        int
	AwardRuleValue string
	End            bool
}

type RuleLogicCheckTypeVO string

type TreeActionEntity struct {
	RuleLogicCheckTypeVO RuleLogicCheckTypeVO
	StrategyAwardVO      StrategyAwardVO
}
type Logic func(userId string, strategyId int64, awardId int, ruleValue string) TreeActionEntity
type Process func()

const (
	ALLOW     RuleLogicCheckTypeVO = "放行；执行后续的流程，不受规则引擎影响"
	TAKE_OVER RuleLogicCheckTypeVO = "接管；后续的流程，受规则引擎执行结果影响"
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

type treeActionEntity struct {
	RuleLogicCheckTypeVO RuleLogicCheckTypeVO
	StrategyAwardVO      StrategyAwardVO
}
