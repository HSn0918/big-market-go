package model

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/hsn0918/BigMarket/common"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StrategyModel = (*customStrategyModel)(nil)

type (
	// StrategyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStrategyModel.
	StrategyModel interface {
		strategyModel
		QueryStrategy(ctx context.Context, StrategyId int64) (strategy *Strategy, err error)
	}

	customStrategyModel struct {
		*defaultStrategyModel
	}
)

func (s Strategy) GetRuleModels() (ruleModels []string) {
	ruleModels = strings.Split(s.RuleModels.String, common.SPLIT)
	return
}
func (s Strategy) GetRuleWeight() (ruleModel string) {
	ruleModels := s.GetRuleModels()
	if ruleModels == nil {
		return ""
	}
	for _, ruleModel := range ruleModels {
		if ruleModel == "rule_weight" {
			return ruleModel
		}
	}
	return ""
}

// NewStrategyModel returns a model for the database table.
func NewStrategyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StrategyModel {
	return &customStrategyModel{
		defaultStrategyModel: newStrategyModel(conn, c, opts...),
	}
}
func (m *customStrategyModel) QueryStrategy(ctx context.Context, StrategyId int64) (strategy *Strategy, err error) {
	strategy = &Strategy{}
	cacheKey := fmt.Sprintf("strategy:%d", StrategyId) // 构造缓存键
	err = m.CachedConn.QueryRowCtx(ctx, strategy, cacheKey, func(context context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := `SELECT * FROM ` + m.table + ` WHERE strategy_id = ? LIMIT 1`
		return conn.QueryRowCtx(ctx, v, query, StrategyId)
	})
	switch {
	case err == nil:
		return strategy, nil // 注意这里返回 strategy
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
