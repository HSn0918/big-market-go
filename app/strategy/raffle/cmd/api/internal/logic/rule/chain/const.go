package chain

type LogicModel string

const (
	RULE_DEFAULT   LogicModel = "默认抽奖"
	RULE_BLACKLIST LogicModel = "黑名单抽奖"
	RULE_WEIGHT    LogicModel = "权重规则"
)

func (m LogicModel) Code() string {
	switch m {
	case RULE_DEFAULT:
		return "rule_default"
	case RULE_BLACKLIST:
		return "rule_blacklist"
	case RULE_WEIGHT:
		return "rule_weight"
	default:
		return "unknown"
	}
}
