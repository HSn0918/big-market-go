package model

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleTreeNodeModel = (*customRuleTreeNodeModel)(nil)

type (
	// RuleTreeNodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleTreeNodeModel.
	RuleTreeNodeModel interface {
		ruleTreeNodeModel
		QueryRuleTreeNodesByTreeId(ctx context.Context, treeId string) (ruleTreeNodes []*RuleTreeNode, err error)
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
func (m *customRuleTreeNodeModel) QueryRuleTreeNodesByTreeId(ctx context.Context, treeId string) (ruleTreeNodes []*RuleTreeNode, err error) {
	ruleTreeNodes = []*RuleTreeNode{}
	query := `SELECT * FROM` + m.table + `WHERE tree_id = ? `
	err = m.QueryRowsNoCacheCtx(ctx, &ruleTreeNodes, query, treeId)
	switch {
	case err == nil:
		return ruleTreeNodes, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, err
	default:
		return nil, err
	}
}
