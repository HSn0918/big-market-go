package model

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleTreeNodeLineModel = (*customRuleTreeNodeLineModel)(nil)

type (
	// RuleTreeNodeLineModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleTreeNodeLineModel.
	RuleTreeNodeLineModel interface {
		ruleTreeNodeLineModel
		QueryRuleTreeNodeLinesByTreeId(ctx context.Context, treeId string) (ruleTreeNodeLines []*RuleTreeNodeLine, err error)
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
func (m *customRuleTreeNodeLineModel) QueryRuleTreeNodeLinesByTreeId(ctx context.Context, treeId string) (ruleTreeNodeLines []*RuleTreeNodeLine, err error) {
	ruleTreeNodeLines = []*RuleTreeNodeLine{}
	query := `select * from` + m.table + `where tree_id = ?`
	err = m.QueryRowsNoCacheCtx(ctx, &ruleTreeNodeLines, query, treeId)
	switch {
	case err == nil:
		return ruleTreeNodeLines, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, err
	default:
		return nil, err
	}
}
