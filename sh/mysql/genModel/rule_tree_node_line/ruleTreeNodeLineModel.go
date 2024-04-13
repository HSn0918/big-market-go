package rule_tree_node_line

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleTreeNodeLineModel = (*customRuleTreeNodeLineModel)(nil)

type (
	// RuleTreeNodeLineModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleTreeNodeLineModel.
	RuleTreeNodeLineModel interface {
		ruleTreeNodeLineModel
	}

	customRuleTreeNodeLineModel struct {
		*defaultRuleTreeNodeLineModel
	}
)

// NewRuleTreeNodeLineModel returns a model for the database table.
func NewRuleTreeNodeLineModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RuleTreeNodeLineModel {
	return &customRuleTreeNodeLineModel{
		defaultRuleTreeNodeLineModel: newRuleTreeNodeLineModel(conn, c, opts...),
	}
}
