package chain

import (
	"context"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

type LogicModel string
type StrategyAwardVO struct {
	AwardId    int
	LogicModel string
	End        bool
}
type LogicChainFunc func(context.Context, *svc.ServiceContext, int64) (StrategyAwardVO, error)

const (
	RULE_DEFAULT   LogicModel = "默认抽奖"
	RULE_BLACKLIST LogicModel = "黑名单抽奖"
	RULE_WEIGHT    LogicModel = "权重规则"
	RULE_ERROR     LogicModel = "规则错误"
)

func CheckStrategyAwardContinue(code string) bool {
	switch code {
	case RULE_BLACKLIST.Code():
		return false
	case RULE_WEIGHT.Code():
		return false
	default:
		return true
	}
}
func (m LogicModel) Code() string {
	switch m {
	case RULE_DEFAULT:
		return "rule_default"
	case RULE_BLACKLIST:
		return "rule_blacklist"
	case RULE_WEIGHT:
		return "rule_weight"
	case RULE_ERROR:
		return "rule_error"

	default:
		return "unknown"
	}
}
