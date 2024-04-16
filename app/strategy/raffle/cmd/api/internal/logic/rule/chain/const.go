package chain

import (
	"context"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

type LogicModel string
type strategyAwardVO struct {
	AwardId    int
	LogicModel string
	End        bool
}
type LogicChainFunc func(context.Context, *svc.ServiceContext, int64) (strategyAwardVO, error)

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
