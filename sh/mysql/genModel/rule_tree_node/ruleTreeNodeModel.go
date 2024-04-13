package rule_tree_node

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleTreeNodeModel = (*customRuleTreeNodeModel)(nil)

type (
	// RuleTreeNodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleTreeNodeModel.
	RuleTreeNodeModel interface {
		ruleTreeNodeModel
	}

	customRuleTreeNodeModel struct {
		*defaultRuleTreeNodeModel
	}
)

// NewRuleTreeNodeModel returns a model for the database table.
func NewRuleTreeNodeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RuleTreeNodeModel {
	return &customRuleTreeNodeModel{
		defaultRuleTreeNodeModel: newRuleTreeNodeModel(conn, c, opts...),
	}
}
