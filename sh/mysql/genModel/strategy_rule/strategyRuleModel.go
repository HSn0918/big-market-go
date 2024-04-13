package strategy_rule

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StrategyRuleModel = (*customStrategyRuleModel)(nil)

type (
	// StrategyRuleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStrategyRuleModel.
	StrategyRuleModel interface {
		strategyRuleModel
	}

	customStrategyRuleModel struct {
		*defaultStrategyRuleModel
	}
)

// NewStrategyRuleModel returns a model for the database table.
func NewStrategyRuleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StrategyRuleModel {
	return &customStrategyRuleModel{
		defaultStrategyRuleModel: newStrategyRuleModel(conn, c, opts...),
	}
}
