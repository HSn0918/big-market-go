package strategy

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StrategyModel = (*customStrategyModel)(nil)

type (
	// StrategyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStrategyModel.
	StrategyModel interface {
		strategyModel
	}

	customStrategyModel struct {
		*defaultStrategyModel
	}
)

// NewStrategyModel returns a model for the database table.
func NewStrategyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StrategyModel {
	return &customStrategyModel{
		defaultStrategyModel: newStrategyModel(conn, c, opts...),
	}
}
