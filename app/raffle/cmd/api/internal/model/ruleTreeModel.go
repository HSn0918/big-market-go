package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleTreeModel = (*customRuleTreeModel)(nil)

type (
	// RuleTreeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleTreeModel.
	RuleTreeModel interface {
		ruleTreeModel
	}

	customRuleTreeModel struct {
		*defaultRuleTreeModel
	}
)

// NewRuleTreeModel returns a model for the database table.
func NewRuleTreeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RuleTreeModel {
	return &customRuleTreeModel{
		defaultRuleTreeModel: newRuleTreeModel(conn, c, opts...),
	}
}
